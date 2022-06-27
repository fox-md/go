package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var Version = "N/A"
var Hash = "N/A"
var BuildDate = "N/A"

type appInfo struct {
	Version   string
	GitHash   string
	BuildDate string
}

type podInfo struct {
	Node      string
	Pod       string
	Namespace string
	PodIp     string
	PodSA     string
}

func getPodInfo() (podInfo podInfo) {
	node := os.Getenv("MY_NODE_NAME")
	if len(node) == 0 {
		node = "N/A"
	}

	pod := os.Getenv("MY_POD_NAME")
	if len(pod) == 0 {
		pod = "N/A"
	}

	ns := os.Getenv("MY_POD_NAMESPACE")
	if len(ns) == 0 {
		ns = "N/A"
	}

	ip := os.Getenv("MY_POD_IP")
	if len(ip) == 0 {
		ip = "N/A"
	}

	sa := os.Getenv("MY_POD_SERVICE_ACCOUNT")
	if len(sa) == 0 {
		sa = "N/A"
	}

	podInfo.Node = node
	podInfo.Pod = pod
	podInfo.Namespace = ns
	podInfo.PodIp = ip
	podInfo.PodSA = sa

	return
}

func getIpAddress() string {
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		log.Print(err.Error())
		return err.Error()
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err.Error())
		}
		sb := string(body) + "\n"
		log.Println(sb)
		return sb
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(getIpAddress()))
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	info := appInfo{
		Version:   Version,
		GitHash:   Hash,
		BuildDate: BuildDate,
	}

	jsonResp, err := json.Marshal(info)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func handlePodInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	podInfo := getPodInfo()

	jsonResp, err := json.Marshal(podInfo)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func main() {

	var enableTLS bool
	var err error

	flag.BoolVar(&enableTLS, "tls", false, "Enable tls")
	flag.Parse()

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/about", handleInfo)
	http.HandleFunc("/pod", handlePodInfo)

	if enableTLS {
		err = http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
		log.Print("Starting TLS listener on port 8443")
	} else {
		err = http.ListenAndServe(":8080", nil)
		log.Print("Starting listener on port 8080")
	}
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
