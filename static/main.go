package main

import (
    "io/ioutil"
    "net/http"
    "log"
)

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
      log.Printf(sb)
      return sb
   }
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(getIpAddress()))
}

func main() {
    http.HandleFunc("/", HelloServer)
    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
