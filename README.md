# GoFunko Scraper

Golang tool based on Colly to build a personnal funko pop toys MySQL database as well as image files.

Dockerized for convenience

## Container Usage

```bash
docker-compose up --build
docker exec -it golang_app bash
```

## Running your go files

### From container

```bash
go run scraper.go
```

### From your computer

```bash
docker exec -it golang_app bash -c "go run scraper.go"
```
