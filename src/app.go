package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	NB_SERVER    = 5
	DEFAULT_FLV  = 1234
	DEFAULT_RTMP = 1935
	DOCKER_IMG   = "eip:server1"
)

var pnsEnv []ServerEnv = make([]ServerEnv, NB_SERVER)
var maxId = int64(0)

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

func launch_crtmpd(id int, quit_chan chan int) {
	// docker run --rm -p 1234:{{.PortFLV}} -p 1935:{{.PortRTMP}} eip:server1
	portflv := putValueInTemplate("{{.PortFLV}}:1234", pnsEnv[id].config)
	portrtmp := putValueInTemplate("{{.PortRTMP}}:1935", pnsEnv[id].config)
	//cmd := exec.Command("docker", "run", "--rm", "--name", pnsEnv[id].config.Name, "-p", portflv, "-p", portrtmp, DOCKER_IMG)
	cmd := exec.Command("docker", "run", "-d", "--name", pnsEnv[id].config.Name, "-p", portflv, "-p", portrtmp, DOCKER_IMG)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal("Can't start the crtmpd server. ", err)
	}
	lol := make([]byte, 4096)
	stderr.Read(lol)
	log.Printf(string(lol))
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	/*
		log.Printf("Waiting for command to finish...")
		go wait_for_app(cmd, quit_chan)
		select {
		case <-quit_chan:
			pnsEnv[id].del()
		}
	*/
}

func wait_for_app(c *exec.Cmd, quit_chan chan int) {
	c.Wait()
	quit_chan <- 1
}

func getPortsForId(id int) (int, int) {
	return DEFAULT_FLV + id, DEFAULT_RTMP + id
}

/*
POST /create
	-> name
	=> { ip, portRTMP, portFLV, id }
*/
func new_server(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	log.Printf("New Server: %s", params["name"])
	// EMPECHER 2x le meme name
	id := getFirstFree()
	if id == -1 {
		w.Write([]byte("{ \"error\": \"1\", \"error_msg\": \"Server is full\" }"))
		return
	}

	portflv, portrtmp := getPortsForId(id)
	pnsEnv[id].Add(params["name"], portflv, portrtmp)
	// VOIR SI Y'A UNE ERREUR OU PAS !
	go launch_crtmpd(id, pnsEnv[id].server_channel)
	w.Write([]byte(strconv.FormatInt(pnsEnv[id].config.id, 10)))
}

func getFirstFree() int {
	for i, val := range pnsEnv {
		if val.used == false {
			return i
		}
	}
	return -1
}

/*
GET /list
	=> json of the containers
*/
func list_server(w http.ResponseWriter, req *http.Request) {
	list := ""
	// MODIFIER CECI POUR ENVOYER DU JSON
	for _, val := range pnsEnv {
		list += fmt.Sprintf("* %v\n", val.config)
	}
	w.Write([]byte(list))
}

func getServerConfigForId(id int64) (*ServerConfig, int64) {
	for id_container, val := range pnsEnv {
		if val.used == true && val.config.id == id {
			return val.config, int64(id_container)
		}
	}
	return nil, int64(-1)
}

/*
POST /delete
	-> id
	=> error in json
*/
func delete_server(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := strconv.ParseInt(params["id"], 10, 32)
	container, id_container := getServerConfigForId(id)

	if container == nil {
		w.Write([]byte("{\"error\" : 2, \"error_msg\": \"Cannot find a container with this id\"}"))
		return
	}

	cmd := exec.Command("docker", "rm", "-f", container.Name)
	err := cmd.Start()
	if err != nil {
		log.Fatal("Can't close the crtmpd server.")
	}
	pnsEnv[id_container].Del() // RECUP ERREUR SI Y'EN A UNE
	log.Printf("deleting id: %d", id)
}

func main() {

	//cmd := exec.Command("docker", "build", "-t", DOCKER_IMG, ".")
	//err := cmd.Run()
	//if err != nil {
	//	log.Fatal("Can't create the Docker image... => ", err)
	//} else {
	//	log.Print("Docker image created.")
	//}

	maxId = 0

	r := mux.NewRouter()
	r.HandleFunc("/new/{name}", new_server).Methods("GET")
	r.HandleFunc("/delete/{id}", delete_server).Methods("GET")
	r.HandleFunc("/list", list_server).Methods("GET")
	log.Printf("Listening ...")
	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}
