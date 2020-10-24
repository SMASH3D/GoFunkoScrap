package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

/*
The Licence struct
*/
type Licence struct {
	Name      string
	Logo      string
	LicenceID int64
	CrawledAt time.Time
	URL       string
}

/*
The Funko struct
*/
type Funko struct {
	Name      string
	ImgURL    string
	Price     float64
	LicenceID int64  //ID  of licence
	Produced  string //date of release
	Scale     string
	Edition   string //translucent...
	Ref       string  //funko ref
	Num       int64  //number within licence
	CrawledAt time.Time
}

//Extracts a numeric ID from a string and a given regular expression
func getIDFromURL(url string, regxpr string) (int64, error) {
	regex := *regexp.MustCompile(regxpr)
	submatches := regex.FindStringSubmatch(url)
	if len(submatches) == 0 {
		return 0, fmt.Errorf("Could not parse id from url : %s", url)
	}
	id, err := strconv.ParseInt(regex.FindStringSubmatch(url)[len(submatches)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse id from url : %s - Reason: %s", url, err)
	}
	return id, err
}

//Builds an array of Licences when parsing a licence list page
func parseLicences() []Licence {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.placedespop.com"),
	)
	licences := make([]Licence, 0)
	//Parsing licences
	c.OnHTML("div.wrapper.wrapper-lics > div.lics > a", func(e *colly.HTMLElement) {
		licence := Licence{}
		url := e.Attr("href")
		licence.Name = e.ChildText(".licl-txt")
		licence.Logo = e.ChildAttr(".licl-logo > img", "src")
		licence.URL = url

		if licenceID, err := getIDFromURL(url, `(?s)\/(\d+)\z`); err != nil {
			fmt.Println("An error occured: ", err)
		} else {
			licence.LicenceID = licenceID
		}

		licence.CrawledAt = time.Now()

		licences = append(licences, licence)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Licences list page : ", r.URL)
	})

	// Start scraping
	c.Visit("https://www.placedespop.com/licences-figurines-funko-pop")

	return licences
}

func scrapeFunkos(licences []Licence) ([]Funko, int) {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.placedespop.com"),
	)
	pageCount := 0
	funkos := make([]Funko, 0)

	c.OnHTML("#TriVoirForm > strong", func(e *colly.HTMLElement) {
		amount := e.Text
		fmt.Println("items : ", amount)
	})

	c.OnHTML("div.wrapper.wrapper-prods > div.prods > a", func(e *colly.HTMLElement) {

		funko := Funko{}
		link := e.Request.URL.String()
		if licenceID, err := getIDFromURL(link, `(?s)\/(\d+)\z`); err != nil {
			fmt.Println("An error occured: ", err)
		} else {
			funko.LicenceID = licenceID
		}

		funko.Name = e.ChildText(".prodl-libelle")
		funko.ImgURL = e.ChildAttr(".prodl-img > img", "data-src")
		if num, err := strconv.ParseInt(strings.ReplaceAll(e.ChildText(".prodl-ref"), "#", ""), 10, 64); err == nil {
			funko.Num = num
		}
		funkoLink := e.Attr("href")
		if funkoRef, err := getIDFromURL(funkoLink, `(\d+)\D+\d*\/\d+$`); err != nil {
			fmt.Println("An error occured: ", err)
			funko.Ref = fmt.Sprintf("UNKNOWN-%d", time.Now().UnixNano() / int64(time.Millisecond))
			fmt.Println("made up Ref: ", funko.Ref)
		} else {
			funko.Ref = strconv.FormatInt(funkoRef, 10)
		}
		if price, err := strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(e.ChildText(".prodl-prix > span"), "€", "")), 64); err == nil {
			funko.Price = price
		}
		funko.CrawledAt = time.Now()
		funkos = append(funkos, funko)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting licence detail page : ", r.URL)
		pageCount++
	})
	for i, licence := range licences {
		// UNCOMMENT FOLLOWING WHEN DEBUGING, TO ONLY PARSE 1 LICENCE PAGE
		//if licence.Name == "Naruto" { 
			fmt.Printf("%d = Licence: %s (%d) \n", i, licence.Name, licence.LicenceID)
			c.Visit(licence.URL)
		//}

	}

	return funkos, pageCount
}

// SaveLicences persists the licence into DB
func SaveLicences(licences []Licence) {
	db, err := sql.Open("mysql", "crawler:popopop@tcp(db:3306)/funkoscrap")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("REPLACE INTO licences(LicenceID, Name, Logo, URL, CrawledAt) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for i, licence := range licences {
		fmt.Printf("%d = Licence: %s (%d) \n", i, licence.Name, licence.LicenceID)
		_, err := stmt.Exec(licence.LicenceID, licence.Name, licence.Logo, licence.URL, licence.CrawledAt.Format(time.RFC3339))
		if err != nil {
			log.Fatal(err)
		}
	}
}

// SaveFunkos persists the licence into DB
func SaveFunkos(funkos []Funko) {
	db, err := sql.Open("mysql", "crawler:popopop@tcp(db:3306)/funkoscrap")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("REPLACE INTO funkos(LicenceID, Ref, Num, Name, ImgURL, Price, CrawledAt) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for i, funko := range funkos {
		fmt.Printf("%d = funko: %s (%s) \n", i, funko.Name, funko.Ref)
		_, err := stmt.Exec(funko.LicenceID, funko.Ref, funko.Num, funko.Name, funko.ImgURL, funko.Price, funko.CrawledAt.Format(time.RFC3339))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	//HANDLING FLAGS
	isVerboseMode := flag.Bool("v", false, "verbose mode")
	flag.Parse()

	licences := parseLicences()
	if *isVerboseMode {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		// Dump json to the standard output
		enc.Encode(licences)
	}
	SaveLicences(licences)
	funkos, pageCount := scrapeFunkos(licences)
	SaveFunkos(funkos)
	if *isVerboseMode {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		// Dump json to the standard output
		//enc.Encode(licences)
		enc.Encode(funkos)
	}

	fmt.Println("pageCount : ", pageCount)
}
