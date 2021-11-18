package main

import (
	"./lcu"
	"fmt"
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
		fmt.Println("Please run this as administrator")
		os.Exit(1)
	} else {
		fmt.Println("Welcome")
	}
}

var lcuInstance = lcu.NewLCUInstance()

const (
	GNone        = "None"
	GLobby       = "Lobby"
	GMatchmaking = "Matchmaking"
	GReadyCheck  = "ReadyCheck"
	GChampSelect = "ChampSelect"
	GInProgress  = "InProgress"
)

func AutoAccept() {
	body, err := lcuInstance.MakeRequest("/lol-gameflow/v1/gameflow-phase")
	if err != nil {
		panic(err)
	}
	if strings.Contains(string(body), GReadyCheck) {
		body, err := lcuInstance.MakePost("/lol-matchmaking/v1/ready-check/accept")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(body))
		fmt.Println("Auto Accepted")
	}
	fmt.Println(string(body))
}

func main() {
	checkRole()
	lcuInstance.GrabToken()
	for {
		AutoAccept()
		time.Sleep(1 * time.Second)
	}

	//notify := make(chan uint32)
	//go lcuInstance.WatchLCU(notify)
	//
	//for {
	//	t:=<-notify
	//	switch t {
	//	case 0:
	//	//	if not running, start a goroutine, set global var running to True
	//	case 2:
	//	//	terminate the goroutine, set global var running to False
	//	default:
	//	//	do nothing
	//	}
	//}
}
