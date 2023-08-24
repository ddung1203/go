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