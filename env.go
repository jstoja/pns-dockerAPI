package main

type ServerEnv struct {
	server_channel chan int
	used           bool
	config         *ServerConfig
}

func (self *ServerEnv) Add(name string, flv int, rtmp int) {
	self.used = true
	self.server_channel = make(chan int)
	maxId += 1
	// MIEUX SELECTIONNER LES PORTS
	self.config = &ServerConfig{maxId, name, rtmp, flv}
}

func (self *ServerEnv) Del() {
	self.used = false
	self.config = nil
}
