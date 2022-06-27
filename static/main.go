package main

import (
        "encoding/json"
        "io/ioutil"
        "log"
        "net/http"
)

var Version = "N/A"
var Hash = "N/A"
var BuildDate = "N/A"

type appInfo struct {
        Version   string
        GitHash   string
        BuildDate string
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

func main() {
        http.HandleFunc("/", HelloServer)
        http.HandleFunc("/about", handleInfo)
        err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}
