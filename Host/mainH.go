package main

//TODO: Fix the many issues! Firstly, figure out a way to let the name (line 50) work. Then, figure out a return thing! Also reduce the ping time and/or make that a setup step. Maybe a Trello board is in order.
import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type AgentFormat struct {
	AgentName string `json:"agentName"`
	AgentAddress string `json:"agentAddress"`	
}
type Responses struct {
	Name string `json:"name"`
	Address string `json:"address"`	
	CPU int `json:"processing"`
	RAM int `json:"memory"`
	Device string `json:"device"`
}

func pollAgents(agents []AgentFormat) {
	var incomplete []AgentFormat
	var responses []Responses
	fetchClient := http.Client{
		Timeout: 10 * time.Second,
	}
	 
for i := 0; i < len(agents); i++ {
	fmt.Println("Fetching:", `http://` + agents[i].AgentAddress + ":4354/stats")
		resp, err := fetchClient.Get(`http://` + agents[i].AgentAddress + ":4354/stats")
		if err != nil {
			incomplete = append(incomplete, agents[i])
			fmt.Println(`"` + agents[i].AgentName + `"`, "failed to respond to the initial ping. Logged as incomplete. Error:", err)
		} else {
			
			responseBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			} else {
				var tempResp Responses
				err = json.Unmarshal(responseBytes, &tempResp)
				responses = append(responses, tempResp)
				if(err != nil) {
					fmt.Println(err)
				} else {
					fmt.Println(tempResp.Name, tempResp.CPU, "%", tempResp.Device)
				}
			}
		}
	}
}

func fetchAllAgents(response http.ResponseWriter, request *http.Request) {
	
	agentsRaw, err := os.ReadFile("Vivus/agentsList/vivusagents.json")
	if(err != nil) {
		if(errors.Is(err, os.ErrNotExist)) {
			fmt.Println("The JSON file does not exist! Creating vivusagents.json now...")
			os.MkdirAll("Vivus/agentsList", 0755)
			writeError := os.WriteFile("Vivus/agentsList/vivusagents.json", []byte("[]"), 0755)
			agentsRaw = []byte("[]")
			if (writeError != nil) {
				fmt.Println("Something went wrong when making the written file.", writeError)
			}
		}
	}

	var agents []AgentFormat
	//This is an AMAZING function!
	err = json.Unmarshal(agentsRaw, &agents)

	if(err != nil) {
		fmt.Println(err)
	}


	
	if len(agents) >= 1 {
		pollAgents(agents);
	}

}

func handleIncoming() {
	
	http.Handle("/fetchAll", http.HandlerFunc(fetchAllAgents)) 
	//Checks if something goes wrong with the ListenAndServe function
	log.Fatal(http.ListenAndServe(":4355", nil))
}

func main() {
	handleIncoming();
}

//Online learning resources used for this file: 
//	https://pkg.go.dev/os#MkdirAll
//	https://tutorialedge.net/golang/parsing-json-with-golang/
//	https://gobyexample.com/directories
//	https://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
//	https://pkg.go.dev/net/http
