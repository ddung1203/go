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

#### Dockerfile

소스 코드, `go.mod`, `index.html` 등 대부분 파일을 Docker image 내 추가했다.

```Dockerfile
FROM bitnami/git:2.41.0 as builder

WORKDIR /app
RUN git clone https://github.com/ddung1203/go.git .

FROM golang:alpine3.18

WORKDIR /go/src/github.com/ddung1203/go
COPY --from=builder /app .
EXPOSE 1323

RUN go build main.go
CMD ["./main"]
```

```bash
 jeonj@ubuntu > ~/go/src/github.com/ddung1203/go > master > docker images
REPOSITORY                          TAG                 IMAGE ID       CREATED          SIZE
ddung1203/go                        latest              c49b03f36737   15 minutes ago   388MB
...
```

#### Dockerfile 이미지 크기 줄이기

상기와 같이, docker image의 alpine 이미지를 사용했음에도 `388MB`의 용량을 사용하였다. Go의 경우 실제 빌드하면 실행파일 한 개만 생성이 된다. 즉, 실행 파일 1개 분량으로 크기를 줄일 수 있다.

> Golang은 컴파일 시 의존성이 모두 한 바이너리 파일에 포함된 채로 컴파일된다. 즉 굳이 OS의 구성요소가 필요하지 않기 때문에 scratch 이미지를 사용해서 오직 바이너리 파일만 포함시키면 된다.

Scratch를 베이스 이미지로 사용하기 위해선, `go build` 명령어를 실행 시 하기와 같이 사용해야 한다.

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o main main.go
```

- `CGO_ENABLED`: cgo를 사용하지 않는다.
  - Scratch 이미지에는 C 바이너리조차 없기 때문에, 반드시 cgo를 비활성화 후 빌드해야 한다.
  - [cgo 참고](https://pkg.go.dev/cmd/cgo)
- `GOOS=linux GOARCH=amd64`: OS와 아키텍처 설정
- `-a`: 모든 의존 패키지를 cgo를 사용하지 않도록 재빌드
- `-ldflags '-s'`: 바이너리를 더 경량화하는 Linker 옵션
  - [-ldflags 참고](https://groups.google.com/g/golang-korea/c/bP3ejliyiqQ/m/igHLKFBfX1gJ?pli=1)


```Dockerfile
FROM golang:1.21.0-bullseye as builder

WORKDIR /go/src/github.com/ddung1203/go
RUN git clone https://github.com/ddung1203/go.git .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o main main.go

FROM scratch
WORKDIR /usr/src/app
COPY --from=builder /go/src/github.com/ddung1203/go/ .

CMD ["./main"]
```

```bash
 jeonj@ubuntu > ~/go/src/github.com/ddung1203/go > master > docker images                
REPOSITORY                          TAG                 IMAGE ID       CREATED              SIZE
ddung1203/go                        latest              56f3a27f18fc   About a minute ago   6.25MB
```