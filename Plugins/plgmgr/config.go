package main

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	PluginName string `json:"plugin_name"`
	Path       string `json:"path"`
	Version    string `json:"version"`
}

func ParseCfg(data []byte) []*Config {
	var c []*Config
	err := json.Unmarshal(data, &c)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return c
}
