package internal

import (
	"fmt"
	"time"
)

type MyTime struct{}

func (t MyTime) TimeHook() {
	fmt.Println(time.Date(3022, 1, 1, 0, 0, 0, 0, &time.Location{}))
}

func (t *MyTime) TimePtrHook() {
	fmt.Println(time.Date(3022, 1, 1, 0, 0, 0, 0, &time.Location{}))
}
