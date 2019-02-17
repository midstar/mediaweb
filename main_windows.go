// Main method for Windows systems
package main

import (
	"os"

	"github.com/kardianos/service"
	"github.com/midstar/llog"
)

var logger service.Logger

type program struct {
	webAPI     *WebAPI
	workingDir string
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	if p.workingDir != "" {
		os.Chdir(p.workingDir)
	}
	p.webAPI = mainCommon()
	go p.run()
	return nil
}
func (p *program) run() {
	p.webAPI.Start()
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	p.webAPI.Stop()
	return nil
}

func main() {
	workingDir := ""
	if len(os.Args) > 1 {
		workingDir = os.Args[1]
	}
	svcConfig := &service.Config{
		Name:        "MediaWEB",
		DisplayName: "MediaWEB",
		Description: "WEB server for photos and videos",
	}

	prg := &program{workingDir: workingDir}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		llog.Panic("%s", err)
	}
	s.Run()
}
