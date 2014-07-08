package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	NB_SERVER    = 5
	DEFAULT_FLV  = 1234
	DEFAULT_RTMP = 1935
	DOCKER_IMG   = "pns:server"
)

type ServerConfig struct {
	Id       int64
	Name     string
	PortRTMP int
	PortFLV  int
}

var pnsEnv []ServerEnv = make([]ServerEnv, NB_SERVER)
var maxId = int64(0)

func launch_crtmpd(id int, quit_chan chan int) {
	portflv := putValueInTemplate("{{.PortFLV}}:1234", pnsEnv[id].config)
	portrtmp := putValueInTemplate("{{.PortRTMP}}:1935", pnsEnv[id].config)
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
		AnswerError(w, -1, "Server is full, please wait or take another one.")
		return
	}

	portflv, portrtmp := getPortsForId(id)
	pnsEnv[id].Add(params["name"], portflv, portrtmp)
	// VOIR SI Y'A UNE ERREUR OU PAS !
	go launch_crtmpd(id, pnsEnv[id].server_channel)
	b, err := json.Marshal(pnsEnv[id].config)
	if err != nil {
		AnswerError(w, -2, "Can't answer in JSON...")
		return
	} else {
		w.Write(b)
	}
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

/*
POST /delete
	-> id
	=> json
*/
func delete_server(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := strconv.ParseInt(params["id"], 10, 32)
	container, id_container := getServerConfigForId(id)

	if container == nil {
		AnswerError(w, -3, "Cannot find a container with this id.")
		return
	}

	cmd := exec.Command("docker", "rm", "-f", container.Name)
	err := cmd.Start()
	if err != nil {
		AnswerError(w, -4, "Can't shutdown the container.")
		return
	}
	pnsEnv[id_container].Del() // RECUP ERREUR SI Y'EN A UNE
	log.Printf("deleting id: %d", id)
	AnswerError(w, 1, "Successfuly deleted.")
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
