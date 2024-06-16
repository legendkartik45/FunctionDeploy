package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type FunctionRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Method   string `json:"method"`
}

var functionName = "user-function"

func main() {
	fs := http.FileServer(http.Dir("ui/public"))
	http.Handle("/", fs)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/invoke", invokeHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req FunctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	funcDir := "functions/user"
	if err := os.MkdirAll(funcDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var funcFile string
	switch req.Language {
	case "go":
		funcFile = filepath.Join(funcDir, "main.go")
	case "python":
		funcFile = filepath.Join(funcDir, "main.py")
	case "javascript":
		funcFile = filepath.Join(funcDir, "main.js")
	default:
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	if err := ioutil.WriteFile(funcFile, []byte(req.Code), 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cmd *exec.Cmd
	switch req.Language {
	case "go":
		cmd = exec.Command("go", "build", "-o", "functions/user/"+functionName, funcFile)
	case "python", "javascript":
		cmd = exec.Command("chmod", "+x", funcFile)
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		http.Error(w, string(output), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Function deployed successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func invokeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cmd := exec.Command("./functions/user/" + functionName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(output)
}
