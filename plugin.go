package main

import (
	"fmt"
	"plugin"
)

type PluginMgrCheckFunc func(path string)

var PluginMgrCheck PluginMgrCheckFunc

func init() {
	pluginSo, err := plugin.Open("Plugins/plgmgr/plgmgr.so")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("[*] PluginMgr start\n")
	hookFunc, err := pluginSo.Lookup("Check")
	if err != nil {
		fmt.Println(err)
		return
	}
	PluginMgrCheck = hookFunc.(func(string))
}
