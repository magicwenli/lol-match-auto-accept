package lcu

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	Start uint32 = iota
	On
	Stop
	Off
)

type LCU struct {
	Client *http.Client
	Url    string
	Port   string
	User   string
	Passwd string
}

func NewLCUInstance() LCU {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	lcu := LCU{}
	lcu.Client = &http.Client{Transport: tr}
	lcu.Url = "https://127.0.0.1:"
	lcu.User = "riot"
	return lcu
}

func (lcu *LCU) GrabToken() {
	cmd := exec.Command("powershell", "$cmdline = Get-WmiObject -Class Win32_Process -Filter \"name='LeagueClientUx.exe'\" | Select-Object -Expand CommandLine\nif($cmdline.length -gt 1){\nif($cmdline -match '--app-port=(\\d*)'){\n$port = $Matches[1]\n}\nif($cmdline -match 'remoting-auth-token=([\\w-]*)'){\n$passwd = $Matches[1]\n}\nreturn $port+':'+$passwd\n}")
	//cmd := exec.Command("Get-WmiObject","-Query \"select * from win32_process where name='LeagueClientUx.exe'\" | Format-List -Property CommandLine")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	lcuConfig := strings.Split(string(out), ":")
	lcu.Port = lcuConfig[0]
	lcu.Passwd = strings.TrimSuffix(lcuConfig[1], "\r\n")
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (lcu *LCU) MakeRequest(endpoint string) ([]byte, error) {
	url := lcu.Url + lcu.Port + endpoint
	req, _ := http.NewRequest("Get", url, nil)
	req.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Basic " + basicAuth(lcu.User, lcu.Passwd)},
	}

	res, err := lcu.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, nil
}

func (lcu *LCU) MakePost(endpoint string) ([]byte, error) {
	url := lcu.Url + lcu.Port + endpoint
	req, _ := http.NewRequest("Post", url, nil)
	req.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Basic " + basicAuth(lcu.User, lcu.Passwd)},
	}
	res, err := lcu.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, nil
}

func (lcu *LCU) GetCurrentSummoner() CurrentSummoner {
	s := CurrentSummoner{}
	body, err := lcu.MakeRequest("/lol-summoner/v1/current-summoner")
	if err != nil {
		panic(err)
	}
	_ = json.Unmarshal(body, &s)
	return s
}

func WatchLCU(notify *chan uint32) {
	event := make(chan uint32)
	go func() {
		for {
			cmd := exec.Command("powershell", "$process = Get-Process -Name 'LeagueClientUx' -ErrorAction SilentlyContinue\n    if($null -ne $process){return \"True\"}else{Write-Host \"False\"}")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			if strings.Contains(string(out), "False") {
				event <- Off
			} else {
				event <- On
			}
			time.Sleep(1 * time.Second)
		}
	}()
	former := Off
	for {
		e := <-event
		if e == On && former == Off {
			log.Print("Start")
			*notify <- Start
		} else if e == Off && former == On {
			log.Print("Stop")
			*notify <- Stop
		} else if e == Off && former == Off {
			log.Print("Game inactive")
		}
		former = e

	}
}
