package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"text/template"
)

type JsonError struct {
	code    int
	message string
}

func AnswerError(w http.ResponseWriter, code int, message string) {
	//b, err := json.Marshal(JsonError{1, "Server is full, please wait or take another one."})
	b, err := json.Marshal(JsonError{code, message})
	if err != nil {
		w.Write([]byte("error"))
	} else {
		w.Write(b)
	}
}

func getPortsForId(id int) (int, int) {
	return DEFAULT_FLV + id, DEFAULT_RTMP + id
}

func getFirstFree() int {
	for i, val := range pnsEnv {
		if val.used == false {
			return i
		}
	}
	return -1
}

func getServerConfigForId(id int64) (*ServerConfig, int64) {
	for id_container, val := range pnsEnv {
		if val.used == true && val.config.Id == id {
			return val.config, int64(id_container)
		}
	}
	return nil, int64(-1)
}

func putValueInTemplate(templ string, obj *ServerConfig) string {
	//tmpl, err := template.New("run docker container").Parse("--name {{.Name}} -p 1234:{{.PortFLV}} -p 1935:{{.PortRTMP}} eip:server1")
	tmpl, err := template.New("run docker container").Parse(templ)
	if err != nil {
		panic(err)
	}
	var docker_cmd bytes.Buffer
	err = tmpl.Execute(&docker_cmd, obj)
	if err != nil {
		panic(err)
	}

	return docker_cmd.String()
}
