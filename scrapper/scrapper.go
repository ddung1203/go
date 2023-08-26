package scrapper

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

// Scrape Jobkorea by a term
func Scrape(term string) {
	var baseURL string = "https://www.jobkorea.co.kr/Search/?stext=" + term
	start := time.Now()
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages :=	getPages(baseURL)

	for i := 1; i <= totalPages; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		// Merging two arrays by ...
		jobs = append(jobs, extractedJobs...)
	}
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
	end := time.Now()
	fmt.Println("Time: ", end.Sub(start))
}

func getPage(page int, baseURL string, mainC chan<- []extractedJob) {
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
	title := CleanString(card.Find(".title").Text())
	location := CleanString(card.Find(".name.dev_view").Text())
	hashtag := CleanString(card.Find(".etc").Text())
	c <- extractedJob{
		id: id, 
		title: title, 
		location: location, 
		hashtag: hashtag,
	}
}

func CleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(baseURL string) int{
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
