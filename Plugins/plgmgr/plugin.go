package main

import (
	"fmt"
	funcmap "myProject/hotpatch/func_map"
	"plugin"
)

type Plugin struct {
	ver  PluginVersion
	name string
	path string
}

func (p *Plugin) Version() PluginVersion {
	return p.ver
}

func (p *Plugin) Init() {
	pluginSo, err := plugin.Open("Plugins/" + p.path + "/" + p.name + ".so")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\n Plugin:%v ver:%v start\n", p.name, p.Version())
	hookFunc, err := pluginSo.Lookup("Hook")
	if err != nil {
		fmt.Println(err)
		return
	}
	hookFunc.(func(map[string]uintptr))(funcmap.FuncMap)
}
