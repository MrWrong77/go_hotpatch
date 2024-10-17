package plgmgr

import (
	"fmt"
	funcmap "myProject/hsq/func_map"
	"plugin"
)

type Plugin struct {
	ver  PluginVersion
	name string
}

func (p *Plugin) Version() PluginVersion {
	return p.ver
}

func (p *Plugin) Init() {
	pluginSo, err := plugin.Open("Plugins/" + p.name + "/" + p.name + ".so")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\n Plugin:%v start\n", p.name)
	hookFunc, err := pluginSo.Lookup("Hook")
	if err != nil {
		fmt.Println(err)
		return
	}
	hookFunc.(func(map[string]uintptr))(funcmap.FuncMap)
}
