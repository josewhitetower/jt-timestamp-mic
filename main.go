package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/timestamp/{timestamp}", handleWithTimestamp)
	r.HandleFunc("/api/timestamp/", handleWithOutTimestamp)
	r.HandleFunc("/api/timestamp", handleWithOutTimestamp)
	var dir string
	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	port := getPort()

	log.Println("Server running in port: " + port)
	log.Fatal(http.ListenAndServe(port, r))

}

func handleWithTimestamp(w http.ResponseWriter, r *http.Request) {

	timestamp := mux.Vars(r)["timestamp"]
	globalError := ""
	type Response struct {
		Unix int64  `json:"unix"`
		UTC  string `json:"utc"`
	}
	isTimeStamp := true
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		globalError = err.Error()
		isTimeStamp = false
	}

	if isTimeStamp {
		tm := time.Unix(i, 0)

		res := Response{i, tm.Format("Mon, Jan 2 2006 15:04:05 MST")}
		json, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			globalError = err.Error()
		}
		w.Write(json)

	} else {
		globalError = ""
		const layout = "2006-01-02"

		time, err := time.Parse(layout, timestamp)
		if err != nil {
			globalError = err.Error()
		}

		res := Response{time.Unix(), time.Format("Mon, Jan 2 2006 15:04:05 MST")}
		json, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			globalError = err.Error()
		}
		if globalError == "" {
			w.Write(json)
		}
	}

	if globalError != "" {
		log.Println(globalError)
		type Error struct {
			Error string `json:"error"`
		}
		res := Error{"Invalid date"}
		json, _ := json.Marshal(res)
		w.Write(json)
	}

}

func handleWithOutTimestamp(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Unix int64  `json:"unix"`
		UTC  string `json:"utc"`
	}

	res := Response{time.Now().Unix(), time.Now().Format("Mon, Jan 2 2006 15:04:05 MST")}
	json, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)

}

// GetPort the Port from the environment so we can run on Heroku
func getPort() string {
	port := os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
