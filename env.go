package main

import (
	"strconv"
	"time"
)

type ServerEnv struct {
	server_channel chan int
	used           bool
	config         *ServerConfig
}

func (self *ServerEnv) Add(name string, flv int, rtmp int) {
	self.used = true
	self.server_channel = make(chan int)
	maxId += 1
	name_timed := append([]byte(name), strconv.FormatInt(time.Now().Unix(), 16)...)
	self.config = &ServerConfig{maxId, string(name_timed), rtmp, flv}
}

func (self *ServerEnv) Del() {
	self.used = false
	self.config = nil
}
