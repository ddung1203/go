# URL Checker

#### Slow URL Checker

하기와 같이 작성 시, 컴파일러가 모르는 에러라는 뜻으로 panic 에러를 발생한다.

문제는, 초기화되지 않은 map에 어떤 값을 넣을 수 없다는 뜻이다.

만약, 비어있는 empty map을 만들고 싶으면 `var results = map[string]string{}`으로 작성한다.

혹은, `make()`를 사용한다. `make()`는 map을 만들어주는 함수로 empty map을 초기화하고 싶을 때 작성한다.

```go
var results map[string]string

results["aaa"] = "bbb"
```


`main.go`
```go
package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("error: Request failed")

func main() {
	results := make(map[string]string)

	urls := []string{
		"https://www.google.com/",
		"https://www.aws.com/",
		"https://www.airbnb.com/",
		"https://www.reddit.com/",
		"https://www.naver.com/",
		"https://www.daum.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		result := "OK"
		err := hitURL(url)
		if err != nil {
			result = "FAILED"
		}
		results[url] = result
	}
	for url, result := range(results) {
		fmt.Println(url, result)
	}
}

func hitURL(url string) (error){
	fmt.Println("Checking: ", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return errRequestFailed
	}
	return nil
}
```

상기 코드는 URL을 순차적으로 체크하는 방식이다. URL을 동시에 체크할 수 있도록 
하기 코드를 참고하자.

#### Gorutines

단지 `go`만 추가함으로써 병렬수행이 가능하다.

하지만, `count()`에 모두 `go`를 추가한다면, 아무런 output 없이 종료된다. Gorutines는 프로그램이 작동하는 동안만 유효하다. 즉, **메인 함수가 실행되는 동안만 실행된다.** 이 경우 메인함수는 첫 번째 gorutines를 실행하고, 두 번째 goroutines를 실행하고, 더이상 진행할 작업이 없기 때문에 끝이 난다. 

메인 함수는 goroutines를 기다려주지 않는다. 만약 parallel한 작업을 실행하고자 한다면, 메인 함수는 다른 goroutines를 기다려주지 않기 때문에 메인 함수가 종료된다. 

`main.go`
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go count("aaa")
	count("bbb")
}

func count(person string) {
	for i := 0; i < 10; i ++ {
		fmt.Println(person, i)
		time.Sleep(time.Second)
	}
}
```

[Slow URL Checker](#slow-url-checker)에서 모든 URL에 goroutine을 적용할 수 있다. 하지만, 메인 함수는 goroutine을 기다려주지 않기 때문에, goroutine으로부터 main function까지 어떻게 커뮤니케이션을 하는지는 Channel을 사용한다.

#### Channel

channel의 타입은 `chan`이고 어떤 타입의 데이터를 주고 받을 것이지를 Go에게 알려줘야 한다.

`channel <- person`: `person`을 `channel`이라는 채널로 보낸다.

`<- channel`: 채널로부터 메시지를 가져온다. Go는 몇 개의 goroutine이 실행 중인지 확인을 하며, 채널의 수가 넘는 메시지를 가져온다면 메시지를 받기 위해서 기다릴 수 없기 때문에 교착 상태(deadlock)가 되어 에러가 나온다.


`main.go`
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	channel := make(chan string)
	people := [2]string{"aaa", "bbb"}
	for _, person := range people {
		go check(person, channel)
	}
	for i := 0; i < len(people); i++ {
		fmt.Println(<- channel)
	}
}

func check(person string, channel chan string) {
	time.Sleep(time.Second * 5)
	channel <- person
}
```

#### Fast URL Checker

`chan<-`의 경우 send only이다.

`main.go`
```go
package main

import (
	"fmt"
	"net/http"
)

type requestResult struct {
	url string
	status string
}

func main() {
	results := make(map[string]string)
	c := make(chan requestResult)

	urls := []string{
		"https://www.google.com/",
		"https://www.aws.com/",
		"https://www.airbnb.com/",
		"https://www.reddit.com/",
		"https://www.naver.com/",
		"https://www.daum.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		go hitURL(url, c)
	}
	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}
	for url, status := range results {
		fmt.Println(url, status)
	}
}

func hitURL(url string, c chan<- requestResult) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		c <- requestResult{url:url, status: "FAILED"}
	} else {
		c <- requestResult{url:url, status: "OK"}
	}
}
```
