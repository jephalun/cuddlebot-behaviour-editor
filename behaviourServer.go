package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type SetPointParams struct {
	Addr      string
	Delay     uint16
	Loop      uint16
	Setpoints []uint16
}

type SleepParams struct {
	Addr []string
}

type BehaviourParams struct {
	Name      string
	Data      string
	Overwrite bool
}

type GestureParams struct {
	Name string
}

var behaviourNameToDataMap map[string]string

var DEFAULT_PATH string

func sleepAll() {
	sleepParams := SleepParams{
		Addr: []string{"purr", "headx", "heady", "ribs", "spine"},
	}

	jsonBytes, err := json.Marshal(sleepParams)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		sendSleepCommand(jsonBytes)
	}
}

func setPointCommand(setPt SetPointParams) {
	jsonBytes, err := json.Marshal(setPt)
	if err != nil {
		fmt.Println("error:", err)
	} else {

		go sendSetPointCommand(jsonBytes)
		time.Sleep(time.Millisecond * time.Duration(setPt.Setpoints[0]))

		sleepParams := SleepParams{
			Addr: []string{setPt.Addr},
		}

		jsonBytes, err = json.Marshal(sleepParams)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			sendSleepCommand(jsonBytes)
		}
	}
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
		log.Println(err)
	} else {
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
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
		log.Println(err)
	} else {
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
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
	var setPt SetPointParams
	err := decoder.Decode(&setPt)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(setPt)

		var setPtParams SetPointParams
		setPtParams.Addr = setPt.Addr
		setPtParams.Delay = setPt.Delay
		setPtParams.Loop = setPt.Loop
		setPtParams.Setpoints = setPt.Setpoints
		setPointCommand(setPtParams)
	}
}

func saveBehaviourParams(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// if origin := req.Header.Get("Origin"); origin != "" {
	// 	rw.Header().Set("Access-Control-Allow-Origin", origin)
	// }

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	var behParams BehaviourParams
	err := decoder.Decode(&behParams)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(behParams)
		log.Println("BehaviourParams to save: ", behParams.Name)
		if _, ok := behaviourNameToDataMap[behParams.Name]; ok && !behParams.Overwrite {
			log.Printf("Behaviour \"%s\" already exists, verifying overwrite", behParams.Name)
			rw.Write([]byte("overwrite?"))
		} else {
			behaviourNameToDataMap[behParams.Name] = behParams.Data

			writeBehavioursToFile(DEFAULT_PATH)
		}
	}
}

func loadBehaviourParams(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// if origin := req.Header.Get("Origin"); origin != "" {
	// 	rw.Header().Set("Access-Control-Allow-Origin", origin)
	// }

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	bytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println(err)
	} else {
		if string(bytes) == "defaults" {
			log.Println("Loading default behaviours")

			behsStr, err := loadBehavioursFromFile(DEFAULT_PATH)

			if err != nil {
				log.Println(err)
			} else {
				rw.Write([]byte(behsStr))
			}
		}
	}
}

func gesture(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// if origin := req.Header.Get("Origin"); origin != "" {
	// 	rw.Header().Set("Access-Control-Allow-Origin", origin)
	// }

	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	var gesture GestureParams
	err := decoder.Decode(&gesture)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("gesture: ", gesture)
	}
}

func writeBehavioursToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	fmt.Fprint(w, "[")
	started := false
	for name, data := range behaviourNameToDataMap {
		if started {
			fmt.Fprint(w, ",")
		}
		started = true
		fmt.Fprint(w, "{\"Name\":\""+name+"\", \"Data\":"+data+"}")
	}
	fmt.Fprint(w, "]")

	return w.Flush()
}

func loadBehavioursFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	allBehsString := ""
	for scanner.Scan() {
		var line = scanner.Text()
		//		log.Println("line: ", line)
		allBehsString += line
	}

	return allBehsString, scanner.Err()
}

func main() {
	DEFAULT_PATH = "./DefaultBehaviours.txt"

	//	loadBehavioursFromFile(DEFAULT_PATH)

	behaviourNameToDataMap = make(map[string]string)

	http.HandleFunc("/gesture", gesture)
	http.HandleFunc("/setpoint", setpoint)
	http.HandleFunc("/savebehaviour", saveBehaviourParams)
	http.HandleFunc("/loadbehaviour", loadBehaviourParams)

	http.Handle("/", http.FileServer(http.Dir("./web")))

	//	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets"))))
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("ListenAndServer err: %s\n", err)
	}
}
