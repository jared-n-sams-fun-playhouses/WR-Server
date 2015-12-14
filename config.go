package main

import (
	"runtime"
	"gopkg.in/ini.v1"
)

/**
 * TODO: maybe switch to toml?
 */

type WRConfig struct {
	root string
	ip string
	port string
}

func (c *WRConfig) Populate() (error) {
	var filepath string
	if runtime.GOOS == "darwin" {
		filepath = "./etc/cfg.osx.ini"
	} else if runtime.GOOS == "windows" {
		// havent tested this, shouldn't work? only because of the path separators
		filepath = ".\etc\cfg.win.ini"
	} else {
		filepath = "./etc/cfg.ini"
	}

	cfg, err := ini.Load(filepath)
	if err != nil {
		return err
	}

	server := cfg.Section("Server")
	
	c.root = server.Key("rootdir").String()
	c.ip   = server.Key("ip").String()
	c.port = server.Key("port").String()
	
	return nil
}
