package glay

import (
	"fmt"
	"github.com/communaute-cimi/glay/utils"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type State int
type Action int

type Application struct {
	Name     string `json:"name"`
	Home     string `json:"home"`
	StartCmd string `json:"startcmd"`
	StopCmd  string `json:"stopcmd"`
	User     string `json:"user"`
}

const (
	UP State = iota
	DOWN
	FAILURE
)

const (
	START Action = iota
	STOP
)

func (app Application) State() (state State, err error) {
	pidpath := fmt.Sprintf("%s/server.pid", app.Home)

	if _, err := os.Stat(pidpath); os.IsNotExist(err) {
		return DOWN, err
	}

	tbpid, err := ioutil.ReadFile(pidpath)

	if err != nil {
		return DOWN, err
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(tbpid)))

	if err != nil {
		return -1, err
	}

	if true == utils.PidIsAlive(int(pid)) {
		return UP, nil
	} else {
		return FAILURE, nil
	}
	return
}

func (app Application) ListenPort() (port int, err error) {
	rx := regexp.MustCompile("http.port=(\\d+)")
	configfpath := fmt.Sprintf("%s/conf/application.conf", app.Home)
	configfile, err := ioutil.ReadFile(configfpath)
	if err != nil {
		return
	}
	// TODO: Mettre 1 au lieu de -1 pour avoir le 1er result ?
	resultrx := rx.FindAllSubmatch(configfile, -1)
	if len(resultrx) != 0 {
		port, err := strconv.Atoi(string(resultrx[0][1]))
		if err != nil {
			return port, err
		}
	}
	return
}

// Plugin Nagios.
func NagiosPlugin(maxFailure int, apps []Application) {
	nbfailure := 0

	msg := ""
	for _, app := range apps {
		_, err := app.State()
		if err != nil {
			nbfailure += 1
			msg += app.Name + " "
		}
	}

	if nbfailure != 0 {
		fmt.Printf(msg + " down")
	} else {
		fmt.Printf("Play! OK - Apps run")
	}

	if nbfailure >= maxFailure {
		os.Exit(2)
	} else if nbfailure <= maxFailure && nbfailure != 0 {
		os.Exit(1)
	}
}
