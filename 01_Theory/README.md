# Theory

목차

1. [Main Package](#main-package)
2. [Packages and Imports](#packages-and-imports)
3. [Variables and Constants](#variables-and-constants)
4. [Functions](#functions)
5. [Multiple Return Values](#multiple-return-values)
6. [Naked Returns](#naked-returns)
7. [Defer](#defer)
8. [Loop](#loop)
9. [Variable Expression](#variable-expression)
10. [Switch](#switch)
11. [Pointers](#pointers)
12. [Array](#array)
13. [Maps](#maps)
14. [Structs](#structs)


## Main Package

만약 자신의 프로젝트를 컴파일 하고 싶다면 패키지의 이름에 대해서는 선택권이 없이 `main.go`이다. 따라서, `main.go`로 이름을 생성했다면, 프로젝트를 컴파일 하고 싶다는 뜻이고 사용할 것이란 뜻이다. 목적에 따라 프로젝트 컴파일이 필요없을 수 있다. 이 경우 `main.go`를 굳이 사용하지 않아도 된다(공유를 위한 라이브러리, 오픈소스 기여 등).

main의 경우 Entry point이기 때문에 컴파일러는 패키지의 이름이 main인 것 부터 찾아낸다. 

`main.go`
```go
package main
```

```bash
jeonj@ubuntu > ~/go/src/github.com/ddung1203/go > go run main.go 
# command-line-arguments
runtime.main_main·f: function main is undeclared in the main package
```

Go는 NodeJS와 Python과는 달리 특정 function을 찾게되는데, 이것이 Go 프로그램의 시작점이 되는 부분이다. 자동적으로 컴파일러는 `main package`와 그 안에 있는 `main` function을 먼저 찾고 실행시킨다. 따라서 Entry point는 `main()` 부터 시작한다.

## Packages and Imports

NodeJS, Python은 무언가를 import 할 때, `import xxx from xxx`와 같이 import를 하며 모듈을 export 한다.

Go의 경우, function을 export 하고 싶다면 function을 대문자로 시작하면 된다.

하기와 같이 `something.go`를 생성 후 `main.go`를 실행한다면, `sayBye`는 private이기 때문에 에러를 확인을 할 수 있다.

따라서, Go에서는 대문자로 작성하는 경우는 다른 패키지로부터 export 된 것이라 볼 수 있다.

`something.go`
```go
package something

import "fmt"

func sayBye() {
	fmt.Println("Bye")
}

func SayHello() {
	fmt.Println("Hello")
}
```

`main.go`
```go
package main

import (
	"fmt"

	"github.com/something"
)

func main() {
	fmt.Println("Hello World")
	something.SayHello()
	something.sayBye()
}
```

## Variables and Constants

Variabel: 변수
Constant: 상수

Go는 type 언어이기 때문에 작성한 값의 타입을 찾아내려는 시도를 한다. 따라서 type이 무엇인지 알려주어야 한다. 

```go
package main

import "fmt"

func main() {
	const name1 string = "Jonogseok"
	fmt.Println(name1)

	var name2 string = "Joongseok"
	name2 = "JOONGSEOK"
	fmt.Println(name2)

	name3 := "Joonseok"
	name3 = "asd"
	fmt.Println(name3)
}
```

또한, 변수를 정의할 때, `func` 안에 `name3 := "Joongseok"`과 같은 방법으로도 정의할 수 있다.

상기와 같이 축약시킨 코드를 사용하면 Go가 적절한 type을 찾아준다. 만약, `func` 밖에서 변수를 정의한다면, 축약시킨 코드를 사용할 수 없다.

## Functions

Go는 표준 라이브러리의 종류가 많다.
  
https://golang.org/pkg/

Function의 예시로 `multiply` 함수를 만들어보자.

하기와 같이 어떤 종류의 Value를 받는지, 어떤 종류의 값을 return 하는지 작성해줘야 한다.

```go
package main

import "fmt"

func multiply(a, b int) int { // multiply(a int, b int) 와 동일
	return a * b
}

func main() {
	fmt.Println(multiply(2, 2))
}
```

## Multiple Return Values

Go는 function이 여러 개의 return 값(multiple value)을 가질 수 있다.

```go
package main

import (
	"fmt"
	"strings"
)

func lenAndUpper(name string) (int, string) {
	return len(name), strings.ToUpper(name)
}

func main() {
	totalLength, upperName := lenAndUpper("Joongseok")
	fmt.Println(totalLength, upperName)
}
```

또한, type 앞에 `.`을 입력해주는 방법으로 원하는 만큼의 arguments를 전달할 수 있다.

```go
package main

import (
	"fmt"
)

func repeatMe(words ...string) {
	fmt.Println(words)	 
}

func main() {
	repeatMe("aaa", "bbb", "ccc", "ddd")
}
```

## Naked Returns

Multiple Return Values에선 'int와 string을 return' 할 것임을 명시했다. 하지만, 이를 return문으로 작성하는 대신, func 줄 내에 작성이 가능하다.

```go
package main

import (
	"fmt"
	"strings"
)

func lenAndUpper(name string) (length int, uppercase string) {
	length = len(name)
	uppercase = strings.ToUpper(name)
	return
}

func main() {
	totalLength, upperName := lenAndUpper("Joongseok")
	fmt.Println(totalLength, upperName)
}
```

상기와 같이 작성함으로서 Go가 자동으로 작동할 수 있도록 한다. `Naked Return`은 return 할 variable을 굳이 명시하지 않아도 된다는 뜻이다. 왜냐하면 어떤 variables를 return할 것인지 `(length int, uppercase string)`와 같이 정의해 두었기 떄문이다.

## Defer

function이 끝났을 때 추가적으로 무엇인가 동작하도록 할 수 있다. `defer`는 function이 값을 return 후 실행된다.

```go
package main

import (
	"fmt"
	"strings"
)

func lenAndUpper(name string) (length int, uppercase string) {
	defer fmt.Println("lenAndUpper func returned")
	length = len(name)
	uppercase = strings.ToUpper(name)
	return
}

func main() {
	totalLength, upperName := lenAndUpper("Joongseok")
	fmt.Println(totalLength, upperName)
}
```

## Loop

오직 `for`만 가능하다. 기본적인 `for`은 하기와 같다.

```go
package main

import "fmt"

func superAdd(numbers ...int) int {
	for i := 0; i < len(numbers); i++ {
		fmt.Println(numbers[i])
	}
	return 1
}

func main() {
	superAdd(1, 2, 3, 4, 5, 6)
}
```

또한 하기와 같이 `range`는 array에 loop를 적용할 수 있도록 한다. `range`는 for에서만 적용이 가능하다.

```go
package main

import "fmt"

func superAdd(numbers ...int) int {
	for number := range numbers {
		fmt.Println(number)
	}
	return 1
}

func main() {
	superAdd(1, 2, 3, 4, 5, 6)
}
```

하지만 결과를 확인하면 예상했던 `1 2 3 4 5 6`이 아닌 `0 1 2 3 4 5`가 출력된다. 이는 number가 아닌, index를 출력한다.

```bash
 jeonj@ubuntu > ~/go/src/github.com/ddung1203/go > go run main.go
0
1
2
3
4
5
```

range는 index와 number를 주기 때문이다. 하기 코드를 실행 시 index와 number가 정상적으로 출력됨을 확인할 수 있다.

```go
package main

import "fmt"

func superAdd(numbers ...int) int {
	for index, number := range numbers {
		fmt.Println(index, number)
	}
	return 1
}

func main() {
	superAdd(1, 2, 3, 4, 5, 6)
}
```

또한, 만약 index를 쓰고 싶지 않은 경우, index 대신 `_`를 작성하여 감출 수 있다.

## Variable Expression

일반적인 if 조건문은 하기와 같다.

```go
package main

import "fmt"

func canIDrink(age int) bool {
	koreanAge := age + 2
	if koreanAge >= 18 {
		return true
	}
	return false
}

func main() {
	fmt.Println(canIDrink(16))
}
```

하지만, 하기의 경우엔 Variable Expression을 통해 작성하였다. 이를 통해 if/else 조건에만 사용하기 위해 variable을 생성했음을 알 수 있다. 

```go
package main

import "fmt"

func canIDrink(age int) bool {
	if koreanAge := age + 2; koreanAge >= 18 {
		return true
	}
	return false
}

func main() {
	fmt.Println(canIDrink(16))
}
```

## Switch

기본적인 switch의 구조는 하기와 같다. 또한 if와 같이 variable expression도 사용 가능하다.

```go
package main

import "fmt"

func canIDrink(age int) bool {
	switch age {
	case 10:
		return false
	case 19:
		return true
	}
	return false
}

func main() {
	fmt.Println(canIDrink(20))
}
```

## Pointers

- 값은 메모리에 저장되며 변수는 일종의 alias이다.
- 메모리 주소를 값으로 가진 변수를 포인터라고 부른다.
- 포인터가 가리키는 값을 가져오는 건 역참조라고 한다.
- 메모리 주소를 직접 대입하거나 포인터 연산을 허용하지 않는다.

```go
package main

import "fmt"

func main() {
	a := 2
	b := &a
	*b = 20
	fmt.Println(a, *b)
}
// 20 20
```

## Array

기본 Array

```go
package main

import "fmt"

func main() {
	names := [5]string{"a", "b", "c"}
	// names[3] = "d"
	names[4] = "e"
	fmt.Println(names)
}
// [a b c  e]
```

Array의 크기에 제한없이 element를 추가하고 싶을 땐 Slice data type을 사용한다.

```go
package main

import "fmt"

func main() {
	names := []string{"a", "b", "c"}
	names =	append(names, "d")
	fmt.Println(names)
}
```

## Maps

기본 map

```go
package main

import "fmt"

func main() {
	me := map[string]string{"name":"Joongseok", "age": "27"}
	for key, value := range me {
		fmt.Println(key, value)
	}
}
```

## Structs

Struct는 map보다 유연한 것이 특징이다. Struct는 어떤 struct인지 정의해주어야 한다. 


```go
package main

import "fmt"

type person struct {
	name string
	age int
	favFood []string
}

func main() {
	favFood := []string{"a", "b"}
	person{name: "Joonseok", age: 27, favFood: favFood}
	fmt.Println(me)
	fmt.Println(me.name)
}
```

`person{"Joonseok", 27, favFood}` 방식으로 structure를 만드는 방법은 어떤 value인지에 대해 명확하지 않기 때문에 `person{name: "Joonseok", age: 27, favFood: favFood}`으로 작성하도록 한다.

