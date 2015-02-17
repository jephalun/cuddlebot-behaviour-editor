package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type SetPointParams struct {
	Addr      string
	Delay     uint16
	Loop      uint16
	Setpoints []uint16
}

func mainView(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
	return
}

func setpoint(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// if origin := req.Header.Get("Origin"); origin != "" {
	// 	rw.Header().Set("Access-Control-Allow-Origin", origin)
	// }

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	log.Println(req.Body)

	decoder := json.NewDecoder(req.Body)
	var t SetPointParams
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
	}
	log.Println(t)
}

func main() {
	http.HandleFunc("/setpoint", setpoint)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	//	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets"))))
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("ListenAndServer err: %s\n", err)
	}
}
