package main

import (
	"fmt"
	"myProject/hotpatch/bussiniess"
	"myProject/hotpatch/hook"
	"reflect"
)

//go:noinline
func NewA() {
	fmt.Println("new A")
}

type Func func()

//go:generate go build --gcflags=all=-N  -buildmode=plugin  -o  p1.so  p1.go
func Hook(m map[string]uintptr) {
	_ = bussiniess.A
	bussiniess.GGG = 333
	if hook.AreFunctionTypesEqual(NewA, bussiniess.A) {
		fmt.Println("same")
	}
	ptr, ok := m["myProject/hotpatch/bussiniess.A"]
	if ok {
		origin := hook.ReplaceFunction(ptr, (uintptr)(hook.GetPtr(reflect.ValueOf(NewA))))
		// hook.RrecoverFunction(ptr, origin)
		_ = origin
	}
	ptr2, ok := m["myProject/hotpatch/bussiniess.MyTime.Time"]
	if ok {
		hook.ReplaceFunction(ptr2, (uintptr)(hook.GetPtr(reflect.ValueOf(bussiniess.MyTime.TimeHook))))
		// hook.ReplaceFunction(ptr2, (uintptr)(hook.GetPtr(reflect.ValueOf(internal.MyTime.TimeHook))))
	}

	ptr3, ok := m["myProject/hotpatch/bussiniess.(*MyTime).TimePtr"]
	if ok {
		hook.ReplaceFunction(ptr3, (uintptr)(hook.GetPtr(reflect.ValueOf((*bussiniess.MyTime).TimePtrHook))))
		// hook.ReplaceFunction(ptr3, (uintptr)(hook.GetPtr(reflect.ValueOf((*internal.MyTime).TimePtrHook))))
	}

	// // 1.function
	// err := gohook.Hook(bussiniess.A, NewA, nil)
	// if err != nil {
	// 	fmt.Printf("hook err %v\n", err)
	// }

	// // 2.value recevier method
	// t := bussiniess.MyTime{}
	// err = gohook.HookMethod(t, "Time", internal.MyTime.TimeHook, nil)
	// if err != nil {
	// 	fmt.Printf("hook err %v\n", err)
	// }

	// // 3.pointer recevier method
	// err = gohook.HookMethod(&t, "TimePtr", (*internal.MyTime).TimePtrHook, nil)
	// if err != nil {
	// 	fmt.Printf("hook err %v\n", err)
	// }
}
