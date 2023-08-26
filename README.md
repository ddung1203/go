# Go

Go는 NodeJS나 Python 처럼 원하는 곳에 디렉토리 내 프로젝트를 만들어 사용할 수 없다. 예를 들어, NodeJS는 `package.json`을 가지고서, 우리가 원하는 폴더에 package들을 다운로드 받을 수 있다. 하지만 Go의 경우, `/usr/local/go` 내 Go 트리를 생성하여 작동한다. 따라서 Go의 코드는 Go PATH 디렉토리에 저장되어야 한다.

```bash
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
$ export PATH=$PATH:/usr/local/go/bin

$ go version
go version go1.21.0 linux/amd64
```

Go 폴더를 확인해보면, `src` 폴더 내 도메인을 확인할 수 있다. NodeJS, Python의 경우 npm, pypi로 모듈이나 패키지를 다운로드 받을 수 있는 곳이 한정적인것에 반해, 원하는 곳 어디에서든 코드를 다운로드 받아 사용할 수 있다. Go에서 받아온 코드들을 보기 좋게 정리하는 방법은 도메인별로 분류를 해서 저장하는 것이다.

하기와 같이 `src` 폴더 내 `github.com` 디렉토리를 생성 후, GitHub username으로 폴더를 생성하였다.

```bash
jeonj@ubuntu > ~/go/src > mkdir github.com && cd github.com
jeonj@ubuntu > ~/go/src/github.com > mkdir -p ddung1203/go
```

목차

1. [Theory](./theory/README.md)
2. [Banking](./accounts/README.md)
3. [Dictionary](./mydict/README.md)
4. [URL Checker](./urlchecker/README.md)
5. [Scrapper - Goroutines](./scrapper/README.md)
6. [ECHO - Web Server](./echo/README.md)


#### Go Modules 사용

- `go mod init [module-name]`
  - 명령어의 이름과 같이, Module을 처음 사용할 때 사용한다.
  - `module-name`은 보통 `github.com/[USERNAME]/[REPONAME]` 포맷을 따른다.
  - 본 프로젝트의 경우, 명령어를 실행하여 현재 디렉터리를 Module의 루트로 만든다.

- `go get [module-name]`
  - `go get`을 사용하여 패키지 다운로드

- `go mod tidy`
  - 소스 코드를 확인해서 import 되지 않는 Module들을 자동으로 `go.mod` 파일에서 삭제하고 import 되었지만 실제 Module이 다운안된 경우는 `go.mod`파일에 추가한다.

- `go mod vendor`
  - Module을 이용하면 Module들을 project 밑에 저장하지 않고, GOPATH에 저장하게 된다. 그러나 자신이 이용하던 Module들(자동으로 변경될 수 있는 Module들을 고정시키고 싶을때와 같이)을 repo에 넣고 싶을 경우가 있다. 따라서 이 명령어를 실행시키면 Module들을 자신의 repo 아래 vendor 폴더에 복사를 하게 된다.

```bash
 jeonj@ubuntu > ~/go/src/github.com/ddung1203/go > master > go mod init                          
go: creating new go.mod: module github.com/ddung1203/go
go: to add module requirements and sums:
        go mod tidy

 jeonj@ubuntu > ~/go/src/github.com/ddung1203/go > master > go get github.com/PuerkitoBio/goquery
go: added github.com/PuerkitoBio/goquery v1.8.1
go: added github.com/andybalholm/cascadia v1.3.1
go: added golang.org/x/net v0.7.0
```