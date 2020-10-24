# GoFunko Scraper

Golang tool based on Colly to build a personnal funko pop toys MySQL database as well as image files.

Dockerized for convenience. Such hype.

## Container Usage

```bash
docker-compose up --build
docker exec -it golang_app bash
```

## Installation

When building the containers for the first time, the db needs to be initialized. For some reason it doesn't happen automatically while it should since the scripts are copied into docker-entrypoint-initdb.d

```bash
docker exec -it golang_db bash
mysql -ucrawler  funkoscrap  -ppopopop -e < /sql_files/install.sql
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

From golang_db container
```bash
mysql -ucrawler  funkoscrap  -ppopopop -e "select ImgURL from funkos" | sed 's/\t/","/g;s/^/"/;s/$/"/;s/\n//g' > funkos_images.csv
```

You can still access the db from your host computer using:

```bash
mysql --host=127.0.0.1 -ucrawler  funkoscrap  -ppopopop --port=3307
```