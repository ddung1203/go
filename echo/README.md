# ECHO

https://echo.labstack.com/

Go echo로 서버를 구축한다. 각종 단어도 검색할 수 있도록 서버를 띄워서 검색 기능을 구현하도록 한다.

우선, 기존 `main.go`에 개발한 코드를 가시성을 위해 `/scrapper/scrapper.go`로 옮기고, 새로운 `main.go`에서 import 및 서버를 사용할 수 있도록 한다. `/scrapper/scrapper.go`는 `Scrape`, `CleanString`을 public으로 변경하였으며, `baseURL`을 전역변수에서 지역변수로 변경하는 과정을 거쳤다.

[코드 확인](./../scrapper/scrapper.go)

`main.go`
```go
package main

import (
	"os"
	"strings"

	"github.com/ddung1203/go/scrapper"
	"github.com/labstack/echo"
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error{
	return c.File("home.html")
}

func handleScrape(c echo.Context) error{
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(fileName, fileName)
}

func main(){
	// scrapper.Scrape("kubernetes")
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}
```

