package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type SetPointWeb struct {
	Addr      string
	Delay     uint16 `json:",string"`
	Loop      uint16 `json:",string"`
	Setpoints []uint16
}

type SetPointParams struct {
	Addr      string
	Delay     uint16
	Loop      uint16
	Setpoints []uint16
}

type SleepParams struct {
	Addr []string
}

func sleepAll() {
	sleepParams := SleepParams{
		Addr: []string{"purr", "headx", "heady", "ribs", "spine"},
	}

	jsonBytes, err := json.Marshal(sleepParams)
	if err != nil {
		fmt.Println("error:", err)
	}

	sendSleepCommand(jsonBytes)
}

func setPointCommand(setPt SetPointParams) {
	jsonBytes, err := json.Marshal(setPt)
	if err != nil {
		fmt.Println("error:", err)
	}

	go sendSetPointCommand(jsonBytes)
	time.Sleep(time.Millisecond * time.Duration(setPt.Setpoints[0]))

	sleepParams := SleepParams{
		Addr: []string{setPt.Addr},
	}

	jsonBytes, err = json.Marshal(sleepParams)
	if err != nil {
		fmt.Println("error:", err)
	}

	sendSleepCommand(jsonBytes)
}

func sendSetPointCommand(commandBytes []byte) {
	url := "http://10.10.10.1/1/setpoint.json"

	//  var jsonStr = []byte("{'addr':'purr', 'delay':0, 'loop':65535, 'setpoints':[1000, 28672]}")
	//    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(commandBytes))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func sendSleepCommand(commandBytes []byte) {
	url := "http://10.10.10.1/1/sleep.json"

	//  var jsonStr = []byte("{'addr':'purr', 'delay':0, 'loop':65535, 'setpoints':[1000, 28672]}")
	//    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(commandBytes))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
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

	decoder := json.NewDecoder(req.Body)
	var setPt SetPointWeb
	err := decoder.Decode(&setPt)
	if err != nil {
		log.Println(err)
	}
	log.Println(setPt)

	var setPtParams SetPointParams
	setPtParams.Addr = setPt.Addr
	setPtParams.Delay = setPt.Delay
	setPtParams.Loop = setPt.Loop
	setPtParams.Setpoints = setPt.Setpoints
	setPointCommand(setPtParams)
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
