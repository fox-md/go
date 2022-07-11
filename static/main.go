package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
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
	Version   string
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
	podInfo.Version = Version
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
	w.WriteHeader(http.StatusCreated) // set response code 201
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
	w.WriteHeader(http.StatusOK) // set response code 200
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
	var enableShutdownDelay bool
	var shutdownDelay int

	flag.BoolVar(&enableTLS, "tls", false, "Enable tls")
	flag.BoolVar(&enableShutdownDelay, "delay", false, "Enable Shutdown Delay")
	flag.IntVar(&shutdownDelay, "timeout", 30, "Shutdown delay. Default: 30 sec")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", HelloServer).Methods("GET")
	router.HandleFunc("/about", handleInfo).Methods("GET")
	router.HandleFunc("/pod", handlePodInfo).Methods("GET")

	server := &http.Server{
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if enableTLS {
			log.Print("Starting TLS listener on port 8443")
			server.Addr = ":8443"
			if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		} else {
			log.Print("Starting listener on port 8080")
			server.Addr = ":8080"
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()

	log.Print("Server Started")
	s := <-stop

	switch s {

	case syscall.SIGHUP:
		log.Print("Signal hang up triggered.")

	case syscall.SIGINT:
		log.Print("Signal interrupt triggered.")

	case syscall.SIGTERM:
		log.Print("Signal terminte triggered.")

	case syscall.SIGQUIT:
		log.Print("Signal quit triggered.")

	default:
		log.Print("Unhandel signal.")
	}

	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(shutdownDelay)+5)*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if enableShutdownDelay {
		log.Print("Delay for ", shutdownDelay, " seconds before socket shutdown")
		time.Sleep(time.Duration(shutdownDelay) * time.Second)
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Print("Server Exited Properly")
}
