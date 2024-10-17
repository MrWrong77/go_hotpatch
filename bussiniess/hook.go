package bussiniess

import (
	"fmt"
	"time"
)

//go:noinline
func (t MyTime) TimeHook() {
	fmt.Println(time.Date(2022, 1, 1, 0, 0, 0, 0, &time.Location{}))
}

//go:noinline
func (t *MyTime) TimePtrHook() {
	fmt.Println(time.Date(2033, 1, 1, 0, 0, 0, 0, &time.Location{}))
}
