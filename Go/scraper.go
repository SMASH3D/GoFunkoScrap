package main

import (
  "encoding/json"
  "os"
	"fmt"
	"regexp"
  "strconv"
  "errors"
	"time"
	"github.com/gocolly/colly"
  "flag"
)

type Licence struct {
	Name     	string
	Logo	  	string
	Id	     	int64
	CrawledAt string
	Url       string
}

type Funko struct {
	Name     	string
	ImgURL  	string
	Brand   	string
	Price   	float64
	LicenceId int64
	Produced 	string
	Scale   	string
	Edition  	string
	Ref	     	string
	CrawledAt string
}

//Extracts a numeric ID from a string and a given regular expression
func getIdFromUrl(url string, regxpr string) (int64, error) {
	regex := *regexp.MustCompile(regxpr)
	id, err := strconv.ParseInt(regex.FindStringSubmatch(url)[1], 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Could not parse id from url : %s", url))
	} else {
		return id, nil
	}
}

//Builds an array of Licences when parsing a licence list page
func parseLicences() ([]Licence) {
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
    licence.Url  = url
		// os.MkdirAll(outputDir, os.ModePerm)
    //
		// c.OnResponse(func(r *colly.Response) {
		// 	if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
		// 		r.Save(outputDir + r.FileName())
		// 		return
		// 	}
		// })

		if licenceId, err := getIdFromUrl(url, `(?s)\/(\d+)\z`) ; err != nil {
			fmt.Println("An error occured: ", err)
		} else {
			licence.Id = licenceId
		}

		licence.CrawledAt = time.Now().Format(time.RFC850)

		licences = append(licences, licence)
	})
  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting Licences list page : ", r.URL)
  })

	// Start scraping
	c.Visit("https://www.placedespop.com/licences-figurines-funko-pop")
	//c.Visit("https://www.placedespop.com/figurines-funko-pop/fantastik-plastik/173")

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
  // c.OnHTML("div.lic-desc > div > p", func(e *colly.HTMLElement) {
  //   desc := e.Text
  //   fmt.Println("desc : ", desc)
  // })



  c.OnHTML("div.wrapper.wrapper-prods > div.prods > a", func(e *colly.HTMLElement) {


    funko := Funko{}
    link := e.Request.URL.String()
    if licenceId, err := getIdFromUrl(link, `(?s)\/(\d+)\z`) ; err != nil {
  		fmt.Println("An error occured: ", err)
  	} else {
  		funko.LicenceId = licenceId
  	}



    funko.Name 	= e.ChildText(".prodl-libelle")
    funko.Ref 	= e.ChildText(".prodl-ref")
    if price, err := strconv.ParseFloat(e.ChildText(".prodl-prix"), 64); err == nil {
      funko.Price = price
    }
    funko.CrawledAt = time.Now().Format(time.RFC850)
    funkos = append(funkos, funko)
  })

  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting licence detail page : ", r.URL)
    pageCount++
  })

  for i, licence := range licences {
    fmt.Printf("%d = Licence: %s (%d) \n", i, licence.Name, licence.Id)
    c.Visit(licence.Url)

    if i == 1 {
      break
    }
  }
  return funkos, pageCount
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
  funkos, pageCount := scrapeFunkos(licences)
  if *isVerboseMode {
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "  ")
    // Dump json to the standard output
    //enc.Encode(licences)
    enc.Encode(funkos)
  }
  fmt.Println("pageCount : ", pageCount)
}
