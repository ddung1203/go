# Scrapper - Goroutines

#### getPages

[goquery](https://github.com/PuerkitoBio/goquery)를 이용해 html을 파싱한다.

```bash
go get github.com/PuerkitoBio/goquery
```

`main.go`
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://www.jobkorea.co.kr/Search/?stext=kubernetes"

func main() {
	totalPages :=	getPages()

	for i := 0; i < totalPages; i++ {
		getPage(i)
	}

	fmt.Println(totalPages)
}

func getPage(page int) {
	pageURL := baseURL + "&Page_No=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageURL)
}

func getPages() int{
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".tplPagination.newVer.wide").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}
```

**request**

만약 User-Agent 등의 헤더를 추가하고 싶다면, 하기와 같이 작성한다.

```go
// NewRequest
req, err := http.NewRequest("GET", url, nil)
if err != nil {
 log.Fatal(err)
}
req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
 log.Fatal(err)
}
defer resp.Body.Close()
```

**Find()**

해당되는 노드를 전부 탐색하며, `*goquery.Selection` 구조체를 리턴한다. 이를 순회하기 위해서 `.Each()`를 사용한다.

**Each()**

`Each()`의 파라미터로 `func(int, *goquery.Selection)`이 요구된다.

`Each()`외에도 `EachWithBreak()`가 존재한다. 특정 조건에서 순회를 멈추고 그 Selection을 전달받고 싶을 때 사용한다. 

#### extractJob

strings 패키지를 사용하여 문자열을 수정하도록 한다.

```go
func Join(a []string, sep string) string
// 문자열 슬라이스에 저장된 문자열을 모두 연결

func Fields(s string) []string
// 공백을 기준으로 문자열을 쪼개어 문자열 슬라이스로 저장

func TrimSpace(s string) string
// 문자열 앞뒤에 오는 공백 문자 제거
```

`main.go`
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	hashtag string
}

var baseURL string = "https://www.jobkorea.co.kr/Search/?stext=kubernetes"

func main() {
	var jobs []extractedJob
	totalPages :=	getPages()

	for i := 1; i <= totalPages; i++ {
		extractedJob := getPage(i)
		jobs = append(jobs, extractedJob...)
	}
	fmt.Println(jobs)
}

func getPage(page int) []extractedJob{
	var jobs []extractedJob
	pageURL := baseURL + "&Page_No=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".list-default .list-post")

	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob{
	id, _ := card.Attr("data-gno")
	title := cleanString(card.Find(".title").Text())
	location := cleanString(card.Find(".name.dev_view").Text())
	hashtag := cleanString(card.Find(".etc").Text())
	return extractedJob{
		id: id, 
		title: title, 
		location: location, 
		hashtag: hashtag,
	}
}

func cleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int{
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".tplPagination.newVer.wide").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}
```

#### Writing Jobs

```go
func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"LINK", "TITLE", "LOCATION", "HASHTAG"}
	
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.jobkorea.co.kr/Recruit/GI_Read/"+job.id, job.title, job.location, job.hashtag}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}
```

#### Goroutine

Goroutine을 적용한 코드의 경우, `1.**`의 시간이 소요되는 반면, 순차적으로 읽고 쓰는 코드의 경우 `5.**`의 시간이 소요된다. 

`main.go`
```go
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	hashtag string
}

var baseURL string = "https://www.jobkorea.co.kr/Search/?stext=kubernetes"

func main() {
	start := time.Now()
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages :=	getPages()

	for i := 1; i <= totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		// Merging two arrays by ...
		jobs = append(jobs, extractedJobs...)
	}
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&Page_No=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".list-default .list-post")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-gno")
	title := cleanString(card.Find(".title").Text())
	location := cleanString(card.Find(".name.dev_view").Text())
	hashtag := cleanString(card.Find(".etc").Text())
	c <- extractedJob{
		id: id, 
		title: title, 
		location: location, 
		hashtag: hashtag,
	}
}

func cleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int{
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".tplPagination.newVer.wide").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"LINK", "TITLE", "LOCATION", "HASHTAG"}
	
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.jobkorea.co.kr/Recruit/GI_Read/"+job.id, job.title, job.location, job.hashtag}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}
```