package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/magicwenli/lol-match-auto-accept/lcu"
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
		log.Println()
		log.Print("Press the Enter Key to Quit")
		fmt.Scanln()
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

func AutoAccept(inTL *walk.TextLabel) {
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
			_ = inTL.SetText(string(body))

			if strings.Contains(string(body), GReadyCheck) {
				_, err := lcuInstance.MakePost("/lol-matchmaking/v1/ready-check/accept")
				if err != nil {
					panic(err)
				}
				//log.Print(string(body))
				log.Print("Auto Accepted")
			} else if strings.Contains(string(body), GChampSelect) {
				log.Print("Good for ChampSelect, let me sleep for 1 minute")
				time.Sleep(1 * time.Minute)
			} else if strings.Contains(string(body), GInProgress) {
				log.Print("It seems like game has started, let me sleep for 1 minutes")
				time.Sleep(1 * time.Minute)
			}
			//log.Print(string(body))
			time.Sleep(1 * time.Second)
		}
	}
}

var (
	autoAccepted = false
)

func main() {
	checkRole()

	var mw *walk.MainWindow
	var inTL *walk.TextLabel
	var stCB *walk.CheckBox
	var goCB *walk.CheckBox

	go func() {
		if err := (MainWindow{
			AssignTo: &mw,
			Title:    "LMAA",
			MinSize:  Size{200, 150},
			Size:     Size{250, 200},
			Layout:   VBox{},
			Children: []Widget{
				HSplitter{
					Children: []Widget{
						GroupBox{
							Title:  "Control",
							Layout: Grid{Columns: 1},
							Children: []Widget{
								CheckBox{
									AssignTo:       &goCB,
									Name:           "Game On",
									Text:           "Game Stat",
									TextOnLeftSide: true,
									Checked:        false,
									Enabled:        false,
									Accessibility: Accessibility{
										Help: "Check if Game is Running",
									},
								},
								CheckBox{
									AssignTo:       &stCB,
									Name:           "Game On",
									Text:           "Auto Accept",
									TextOnLeftSide: true,
									Checked:        false,
									Enabled:        false,
									Accessibility: Accessibility{
										Help: "Set Auto Accept",
									},
									OnCheckedChanged: func() {
										if stCB.Checked() { // run
											go AutoAccept(inTL)
											autoAccepted = true
										} else {
											quitCh <- true
											autoAccepted = false
										}
									},
								},
							},
						},
						GroupBox{
							Title:  "Gameflow",
							Layout: Grid{Columns: 1},
							Children: []Widget{
								TextLabel{AssignTo: &inTL},
							},
						},
					},
				},
			},
		}.Create()); err != nil {
			log.Fatal(err)
		}

		lv, err := NewLogView(mw)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(lv)

		icon, _ := walk.Resources.Icon("./assets/icon.ico")
		mw.SetIcon(icon)

		mw.Run()
		os.Exit(0)
	}()

	notify := make(chan uint32)
	go lcu.WatchLCU(&notify)

	for {
		t := <-notify
		switch t {
		case 0:
			if !goCB.Checked() {
				goCB.SetChecked(true)
				stCB.SetEnabled(true)
				log.Print("start goroutine")
			}
		//	if not running, start a goroutine, set global var gameRunning to True
		case 2:
			if goCB.Checked() {
				quitCh <- true
				goCB.SetChecked(false)
				stCB.SetEnabled(false)
				stCB.SetChecked(false)
				log.Print("terminate goroutine")
			}
		//	terminate the goroutine, set global var running to False
		default:
			log.Print("hanging around")
		}
	}
}
