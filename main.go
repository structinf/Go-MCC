// Copyright 2017 Andrew Goulas
// https://www.structinf.com
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var DefaultConfig = &Config{
	Port:       25565,
	Name:       "Go-MCC",
	MOTD:       "Welcome!",
	Verify:     false,
	Public:     true,
	MaxPlayers: 32,
	Heartbeat:  "http://www.classicube.net/heartbeat.jsp",
	MainLevel:  "main",
}

func ReadConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		file, err = json.MarshalIndent(DefaultConfig, "", "\t")
		if err != nil {
			fmt.Printf("Config Error: %s\n", err.Error())
			return DefaultConfig
		}

		err = ioutil.WriteFile(path, file, 0644)
		if err != nil {
			fmt.Printf("Config Error: %s\n", err.Error())
		}

		return DefaultConfig
	} else {
		config := &Config{}
		err = json.Unmarshal(file, config)
		if err != nil {
			fmt.Printf("Config Error: %s\n", err.Error())
			config = DefaultConfig
		}

		return config
	}
}

func main() {
	config := ReadConfig("server.properties")
	storage := NewLvlStorage("levels/")
	server := NewServer(config, storage)
	if server == nil {
		return
	}

	wg := &sync.WaitGroup{}
	go server.Run(wg)

	console := NewConsole(server, wg)
	console.Run()
}
