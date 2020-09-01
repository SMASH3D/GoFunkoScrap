# GoFunko Scraper

Golang tool based on Colly to build a personnal funko pop toys MySQL database as well as image files.

Dockerized for convenience. Such hype.

## Container Usage

```bash
docker-compose up --build
docker exec -it golang_app bash
```

## Running the scraper

### From container

```bash
go run scraper.go
```

### From your computer

```bash
docker exec -it golang_app bash -c "go run scraper.go"
```

## Extracting some data

### Exporting the funko image names into a csv file

```bash
mysql -ucrawler  funkoscrap  -ppopopop -e "select ImgURL from funkos" | sed 's/\t/","/g;s/^/"/;s/$/"/;s/\n//g' > funkos_images.csv
```
