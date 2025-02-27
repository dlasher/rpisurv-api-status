package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "os"
    "strings"

    "github.com/bendahl/uinput"
)

// Response defines the JSON response structure for camera actions.
type Response struct {
    Status string `json:"status"`
    Camera string `json:"camera"`
}

// ErrorResponse defines the JSON error response structure.
type ErrorResponse struct {
    Error string `json:"error"`
}

// allowedNet holds the allowed IP network, initialized at startup.
var allowedNet *net.IPNet

// keyboardDev is the uinput keyboard device.
var keyboardDev uinput.Keyboard

// checkUinputPermissions verifies that the program has read-write access to /dev/uinput.
func checkUinputPermissions() error {
    file, err := os.OpenFile("/dev/uinput", os.O_RDWR, 0)
    if err != nil {
        if os.IsPermission(err) {
            return fmt.Errorf("insufficient permissions to access /dev/uinput. Try running with sudo or adjust permissions with 'sudo chmod 666 /dev/uinput'")
        }
        return fmt.Errorf("failed to open /dev/uinput: %v", err)
    }
    file.Close()
    return nil
}

func init() {
    var err error
    // Parse the allowed IP network (e.g., 10.4.0.0/16)
    allowedNet, err = net.ParseCIDR("10.4.0.0/16")
    if err != nil {
        log.Fatal("Invalid CIDR:", err)
    }

    // Check permissions for /dev/uinput
    if err = checkUinputPermissions(); err != nil {
        log.Fatal(err)
    }

    // Initialize the uinput virtual keyboard device
    keyboardDev, err = uinput.CreateKeyboard("/dev/uinput", []byte("VirtualKeyboard"))
    if err != nil {
        log.Fatal("Failed to create uinput keyboard device:", err)
    }
}

// ipFilter restricts access to the allowed IP network.
func ipFilter(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ipStr := r.RemoteAddr
        ip, _, err := net.SplitHostPort(ipStr)
        if err != nil {
            http.Error(w, "Invalid remote address", http.StatusBadRequest)
            return
        }
        remoteIP := net.ParseIP(ip)
        if remoteIP == nil {
            http.Error(w, "Invalid IP address", http.StatusBadRequest)
            return
        }
        if ip4 := remoteIP.To4(); ip4 != nil {
            remoteIP = ip4
        }
        if !allowedNet.Contains(remoteIP) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next(w, r)
    }
}

// writeJSONError writes a JSON error response.
func writeJSONError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    errResp := ErrorResponse{Error: message}
    json.NewEncoder(w).Encode(errResp)
}

// readFile reads and trims a file’s content.
func readFile(path string) (string, error) {
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(content)), nil
}

// simulateKeyPress simulates a key press using uinput.
func simulateKeyPress(key string) error {
    var keyCode int
    switch key {
    case "p":
        keyCode = uinput.KeyP
    case "r":
        keyCode = uinput.KeyR
    case "F1":
        keyCode = uinput.KeyF1
    case "F2":
        keyCode = uinput.KeyF2
    case "F3":
        keyCode = uinput.KeyF3
    case "F4":
        keyCode = uinput.KeyF4
    default:
        return fmt.Errorf("unsupported key: %s", key)
    }

    // Simulate pressing and releasing the key
    if err := keyboardDev.KeyPress(keyCode); err != nil {
        return err
    }
    return nil
}

// indexHandler handles requests to the root endpoint.
func indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    fmt.Fprint(w, "Server Works!")
}

// cameraHandler handles requests to /camera/<action>.
func cameraHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    action := strings.TrimPrefix(r.URL.Path, "/camera/")
    if action == "" {
        writeJSONError(w, "Missing action", http.StatusBadRequest)
        return
    }

    // Validate the action
    validActions := map[string]bool{
        "pause":  true,
        "resume": true,
        "status": true,
        "1":      true,
        "2":      true,
        "3":      true,
        "4":      true,
    }
    if !validActions[action] {
        writeJSONError(w, "Invalid action", http.StatusBadRequest)
        return
    }

    // Handle actions that write files and simulate key presses
    if action == "pause" || action == "resume" {
        err := ioutil.WriteFile("/tmp/rpi.status", []byte(action), 0644)
        if err != nil {
            log.Printf("Error writing to /tmp/rpi.status: %v", err)
            writeJSONError(w, "File write error", http.StatusInternalServerError)
            return
        }
        key := "p"
        if action == "resume" {
            key = "r"
        }
        if err := simulateKeyPress(key); err != nil {
            log.Printf("Error simulating key press '%s': %v", key, err)
        }
    } else if action == "1" || action == "2" || action == "3" || action == "4" {
        err := ioutil.WriteFile("/tmp/rpi.camera", []byte(action), 0644)
        if err != nil {
            log.Printf("Error writing to /tmp/rpi.camera: %v", err)
            writeJSONError(w, "File write error", http.StatusInternalServerError)
            return
        }
        key := "F" + action
        if err := simulateKeyPress(key); err != nil {
            log.Printf("Error simulating key press '%s': %v", key, err)
        }
    }

    // Read current status and camera values
    status, err := readFile("/tmp/rpi.status")
    if err != nil {
        log.Printf("Error reading /tmp/rpi.status: %v", err)
        writeJSONError(w, "File read error", http.StatusInternalServerError)
        return
    }
    camera, err := readFile("/tmp/rpi.camera")
    if err != nil {
        log.Printf("Error reading /tmp/rpi.camera: %v", err)
        writeJSONError(w, "File read error", http.StatusInternalServerError)
        return
    }

    // Send JSON response
    resp := Response{Status: status, Camera: camera}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func main() {
    http.HandleFunc("/", ipFilter(indexHandler))
    http.HandleFunc("/camera/", ipFilter(cameraHandler))
    log.Println("Starting server on :5000")
    err := http.ListenAndServe(":5000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

    // Ensure the uinput device is closed on exit
    defer keyboardDev.Close()
}
