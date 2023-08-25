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