package main

import (
	"fmt"
	"io"
	"os"
)

// pluginMgr entry point
//
//go:generate go build --gcflags=all=-N --gcflags=all=-l  -buildmode=plugin  -o  plgmgr.so
func Check(path string) {
	GPlgMgr.CheckPath(path)
}

type IPlugin interface {
	Version() PluginVersion
	Init()
}

type PluginManager struct {
	plugs map[string]IPlugin
}

var GPlgMgr PluginManager

func (m *PluginManager) CheckPath(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer f.Close()

	data, _ := io.ReadAll(f)
	cfgs := ParseCfg(data)
	for _, cfg := range cfgs {
		if cfg.PluginName == "" {
			//empty name
			continue
		}
		ver := ParseVersion(cfg.Version)
		if ver == nil {
			// invalid version cfg
			continue
		}

		plg := &Plugin{
			ver:  *ver,
			name: cfg.PluginName,
			path: cfg.Path,
		}
		newPlugin := m.Register(plg.name, plg)
		if newPlugin != nil {
			newPlugin.Init()
		}
	}
}

func (m *PluginManager) Register(pluginName string, plg IPlugin) IPlugin {
	if m.plugs == nil {
		m.plugs = make(map[string]IPlugin)
	}

	if plg == nil {
		return nil
	}

	p, ok := m.plugs[pluginName]
	if ok {
		if !p.Version().Less(plg.Version()) {
			return nil
		}
	}

	m.plugs[pluginName] = plg
	return plg
}
