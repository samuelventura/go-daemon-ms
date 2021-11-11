package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
)

type program struct {
	done chan bool
	exit chan bool
}

func (p *program) Start(s service.Service) (err error) {
	p.exit = make(chan bool)
	p.done = make(chan bool)
	inter := service.Interactive()
	go func() {
		defer close(p.done)
		entry(inter, p.exit)
	}()
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	select {
	case <-p.done:
	case <-time.After(3 * time.Second):
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	//-service install, uninstall, start, stop, restart
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()
	svcConfig := &service.Config{
		Name:        "GoDaemonMs",
		DisplayName: "GoDaemonMs Service",
		Description: "GoDaemonMs https://github.com/samuelventura/go-daemon-ms",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
