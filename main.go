package main

import (
	"fmt"
	"myProject/hotpatch/bussiniess"
	_ "myProject/hotpatch/func_map"
	"time"
)

const PluginCfgFile = "./Plugins/plugin_version.json"

//go:generate go build --gcflags=all=-N -o  main  .
func main() {
	t := time.NewTicker(time.Second)
	hookF := time.After(2 * time.Second)
	mytime := &bussiniess.MyTime{}

	for {
		select {
		case <-t.C:
			bussiniess.A()
			mytime.Time()
			mytime.TimePtr()
			fmt.Println(bussiniess.GGG)
		case <-hookF:
			fmt.Printf("plugin check start ... \n")
			if PluginMgrCheck != nil {
				PluginMgrCheck(PluginCfgFile)
				hookF = time.After(2 * time.Second)
				fmt.Printf("plugin check finish \n")
			}
		}
	}

}
