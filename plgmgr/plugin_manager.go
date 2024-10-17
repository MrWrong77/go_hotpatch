package plgmgr

import (
	"bufio"
	"bytes"
	"os"
)

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
		return
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			break
		}
		const PluginCfgFieldNum = 2
		plgInfo := bytes.Split(line, []byte{'='})
		if len(plgInfo) != PluginCfgFieldNum {
			// plugin cfg format not valid
			continue
		}

		if len(plgInfo[0]) <= 0 {
			//empty name
			continue
		}

		ver := ParseVersion(plgInfo[1])
		if ver == nil {
			// invalid version cfg
			continue
		}

		pluginName := string(plgInfo[0])

		p := &Plugin{ver: *ver, name: pluginName}
		newPlugin := m.Register(p.name, p)
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
