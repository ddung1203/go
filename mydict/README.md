# Dictionary

## map

https://go.dev/blog/maps

type은 method를 가질 수 있기에, `fmt.Println(dictionary["first"])`와 같은 방식이 아닌, method를 사용하여 작성하도록 한다.

상기 blog를 보면 maps에 관해 설명되어 있다. map은 key의 존재여부를 알려주는 방법이 있다. map의 key를 호출하면 value와 boolean으로 존재 여부를 알려준다.

#### search

`mydict/mydict.go`
```go
package mydict

import "errors"

var errNotFound = errors.New("error: Not Found")

// Dictionary type
type Dictionary map[string]string

func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}
```

`main.go`
```go
package main

import (
	"fmt"

	"github.com/ddung1203/go/mydict"
)

func main() {
	dictionary := mydict.Dictionary{"first": "First word"}
	definition, err := dictionary.Search("first")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(definition)
	}
}
```

#### add

`mydict/mydict.go`
```go
// Add a word to the dictionary
func (d Dictionary) Add(word, def string) error{
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}
```

`main.go`
```go
func main() {
	dictionary := mydict.Dictionary{}
	err1 := dictionary.Add("aaa", "aaa")
	if err1 != nil {
		fmt.Println(err1)
	}
	err2 := dictionary.Add("aaa", "aaa")
	if err2 != nil {
		fmt.Println(err2)
	}
	definition, err := dictionary.Search("aaa")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(definition)
	}
}
```

#### update

`mydict/mydict.go`
```go
// Update a word to the dictionary
func (d Dictionary) Update(word, def string) error{
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = def
	case errNotFound:
		return errCantUpdate
	}
	return nil
}
```

`main.go`
```go
func main() {
	dictionary := mydict.Dictionary{}
	err1 := dictionary.Add("aaa", "aaa")
	if err1 != nil {
		fmt.Println(err1)
	}
	err2 := dictionary.Update("aaa", "bbb")
	if err2 != nil {
		fmt.Println(err2)
	}
	definition, err := dictionary.Search("aaa")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(definition)
	}
}
```

#### delete

`mydict/mydict.go`
```go
// Delete a word
func (d Dictionary) Delete(word string) error{
	_, err := d.Search(word)
	switch err {
	case nil:
		delete(d, word)
	case errNotFound:
		return errCantDelete
	}
	return nil
}
```

`main.go`
```go
func main() {
	dictionary := mydict.Dictionary{}
	err1 := dictionary.Add("aaa", "aaa")
	if err1 != nil {
		fmt.Println(err1)
	}
	err2 := dictionary.Update("aaa", "bbb")
	if err2 != nil {
		fmt.Println(err2)
	}
	err3 := dictionary.Delete("aaa")
	if err3 != nil {
		fmt.Println(err3)
	}
	definition, err := dictionary.Search("aaa")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(definition)
	}
}
```

전체 코드는 하기와 같다.

`mydict/mydict.go`
```go
package mydict

import "errors"

var (
	errNotFound = errors.New("error: Not Found")
	errWordExists = errors.New("error: That word already exist")
	errCantUpdate = errors.New("error: Can't update non-existing word")
	errCantDelete = errors.New("error: Can't delete non-existing word")
)

// Dictionary type
type Dictionary map[string]string

func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}

// Add a word to the dictionary
func (d Dictionary) Add(word, def string) error{
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}

// Update a word to the dictionary
func (d Dictionary) Update(word, def string) error{
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = def
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) error{
	_, err := d.Search(word)
	switch err {
	case nil:
		delete(d, word)
	case errNotFound:
		return errCantDelete
	}
	return nil
}
```

`main.go`
```go
package main

import (
	"fmt"

	"github.com/ddung1203/go/mydict"
)

func main() {
	dictionary := mydict.Dictionary{}
	err1 := dictionary.Add("aaa", "aaa")
	if err1 != nil {
		fmt.Println(err1)
	}
	err2 := dictionary.Update("aaa", "bbb")
	if err2 != nil {
		fmt.Println(err2)
	}
	err3 := dictionary.Delete("aaa")
	if err3 != nil {
		fmt.Println(err3)
	}
	definition, err := dictionary.Search("aaa")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(definition)
	}
}
```