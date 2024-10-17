package bussiniess

import (
	"fmt"
	"time"
)

var GGG int = 22

//go:noinline
func A() {
	fmt.Println("old A func")
}

type MyTime struct {
}

//go:noinline
func (t MyTime) Time() {
	fmt.Println(time.Now())
}

func (t *MyTime) TimePtr() {
	fmt.Println(time.Now())
}
