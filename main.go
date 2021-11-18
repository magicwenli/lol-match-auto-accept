package main

import (
	"./lcu"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func checkRole() {
	cmd := exec.Command("powershell", "$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())\n    return $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	if strings.Contains(string(out), "False") {
		log.Print("Please run this as administrator")
		os.Exit(1)
	} else {
		log.Print("Welcome")
	}
}

var (
	lcuInstance = lcu.NewLCUInstance()
	quitCh      = make(chan bool)
)

const (
	GNone        = "None"
	GLobby       = "Lobby"
	GMatchmaking = "Matchmaking"
	GReadyCheck  = "ReadyCheck"
	GChampSelect = "ChampSelect"
	GInProgress  = "InProgress"
)

func AutoAccept() {
	lcuInstance.GrabToken()
	for {
		select {
		case <-quitCh:
			return
		default:
			body, err := lcuInstance.MakeRequest("/lol-gameflow/v1/gameflow-phase")
			if err != nil {
				panic(err)
			}
			if strings.Contains(string(body), GReadyCheck) {
				_, err := lcuInstance.MakePost("/lol-matchmaking/v1/ready-check/accept")
				if err != nil {
					panic(err)
				}
				//log.Print(string(body))
				log.Print("Auto Accepted")
			} else if strings.Contains(string(body), GChampSelect) {
				time.Sleep(10 * time.Minute)
			} else if strings.Contains(string(body), GInProgress) {
				time.Sleep(2 * time.Minute)
			}
			//log.Print(string(body))
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	checkRole()

	notify := make(chan uint32)
	go lcu.WatchLCU(&notify)
	running := false

	for {
		t := <-notify
		switch t {
		case 0:
			if !running {
				go AutoAccept()
				running = true
				log.Print("start goroutine")
			}
		//	if not running, start a goroutine, set global var running to True
		case 2:
			if running {
				quitCh <- true
				running = false
				log.Print("terminate goroutine")
			}
		//	terminate the goroutine, set global var running to False
		default:
			log.Print("hanging around")
		}
	}
}
