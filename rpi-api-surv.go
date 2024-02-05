package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
    "github.com/micmonay/keybd_event"
)

func pressKey(key int) {
    kb, err := keybd_event.NewKeyBonding()
    if err != nil {
        panic(err)
    }

    if keybd_event.GetPlatform() == "linux" {
        time.Sleep(2 * time.Second)
    }

    kb.SetKeys(key)
    err = kb.Launching()
    if err != nil {
        panic(err)
    }

    kb.Clear()
    err = kb.Launching()
    if err != nil {
        panic(err)
    }
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Server Works!")
    })
	
	http.HandleFunc("/camera/status", func(w http.ResponseWriter, r *http.Request) {
		rstatus, err := ioutil.ReadFile("/tmp/rpi.status")
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		rcamera, err := ioutil.ReadFile("/tmp/rpi.camera")
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		fmt.Fprintf(w, `{"status": "%s", "camera": "%s"}`, rstatus, rcamera)
	})

    http.HandleFunc("/camera/pause", func(w http.ResponseWriter, r *http.Request) {
        err := ioutil.WriteFile("/tmp/rpi.status", []byte("pause"), 0644)
        if err != nil {
            fmt.Fprintf(w, "Error")
            return
        }

        rcamera, err := ioutil.ReadFile("/tmp/rpi.camera")
        if err != nil {
            fmt.Fprintf(w, "Error")
            return
        }

        pressKey(keybd_event.VK_P)

        fmt.Fprintf(w, `{"status": "pause", "camera": "%s"}`, rcamera)
    })
    http.HandleFunc("/camera/resume", func(w http.ResponseWriter, r *http.Request) {
        err := ioutil.WriteFile("/tmp/rpi.status", []byte("pause"), 0644)
        if err != nil {
            fmt.Fprintf(w, "Error")
            return
        }

        rcamera, err := ioutil.ReadFile("/tmp/rpi.camera")
        if err != nil {
            fmt.Fprintf(w, "Error")
            return
        }

        pressKey(keybd_event.VK_R)

        fmt.Fprintf(w, `{"status": "resume", "camera": "%s"}`, rcamera)
    })
	http.HandleFunc("/camera/1", func(w http.ResponseWriter, r *http.Request) {
		rstatus, err := ioutil.ReadFile("/tmp/rpi.status")
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		err = ioutil.WriteFile("/tmp/rpi.camera", []byte("1"), 0644)
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		pressKey(keybd_event.VK_F1)
	
		fmt.Fprintf(w, `{"status": "%s", "camera": "1"}`, rstatus)
	})
	http.HandleFunc("/camera/2", func(w http.ResponseWriter, r *http.Request) {
		rstatus, err := ioutil.ReadFile("/tmp/rpi.status")
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		err = ioutil.WriteFile("/tmp/rpi.camera", []byte("1"), 0644)
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		pressKey(keybd_event.VK_F2)
	
		fmt.Fprintf(w, `{"status": "%s", "camera": "2"}`, rstatus)
	})
	http.HandleFunc("/camera/3", func(w http.ResponseWriter, r *http.Request) {
		rstatus, err := ioutil.ReadFile("/tmp/rpi.status")
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		err = ioutil.WriteFile("/tmp/rpi.camera", []byte("1"), 0644)
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		pressKey(keybd_event.VK_F3)
	
		fmt.Fprintf(w, `{"status": "%s", "camera": "3"}`, rstatus)
	})
	http.HandleFunc("/camera/4", func(w http.ResponseWriter, r *http.Request) {
		rstatus, err := ioutil.ReadFile("/tmp/rpi.status")
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		err = ioutil.WriteFile("/tmp/rpi.camera", []byte("1"), 0644)
		if err != nil {
			fmt.Fprintf(w, "Error")
			return
		}
	
		pressKey(keybd_event.VK_F4)
	
		fmt.Fprintf(w, `{"status": "%s", "camera": "4"}`, rstatus)
	})

    http.ListenAndServe(":5000", nil)
}
