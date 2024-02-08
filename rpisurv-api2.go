package main

import (
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/yunginnanet/sendkeys"
)

func pressKey(key string) {
    err := sendkeys.Send(key)
    if err != nil {
        panic(err)
    }
}

func handleCamera(w http.ResponseWriter, r *http.Request, camera string, key string) {
    rstatus, err := ioutil.ReadFile("/tmp/rpi.status")
    if err != nil {
        fmt.Fprintf(w, "Error")
        return
    }

    err = ioutil.WriteFile("/tmp/rpi.camera", []byte(camera), 0644)
    if err != nil {
        fmt.Fprintf(w, "Error")
        return
    }

    pressKey(key)

    fmt.Fprintf(w, `{"status": "%s", "camera": "%s"}`, rstatus, camera)
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Server Works!")
    })

    http.HandleFunc("/camera/1", func(w http.ResponseWriter, r *http.Request) {
        handleCamera(w, r, "1", "{F1}")
    })

    http.HandleFunc("/camera/2", func(w http.ResponseWriter, r *http.Request) {
        handleCamera(w, r, "2", "{F2}")
    })

    http.HandleFunc("/camera/3", func(w http.ResponseWriter, r *http.Request) {
        handleCamera(w, r, "3", "{F3}")
    })

    http.HandleFunc("/camera/4", func(w http.ResponseWriter, r *http.Request) {
        handleCamera(w, r, "4", "{F4}")
    })

    http.ListenAndServe(":5000", nil)
}