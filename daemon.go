package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/samuelventura/go-tools"
)

func Run(dro *DaemonDro, exit chan bool) chan bool {
	traceRecover := func() {
		r := recover()
		if r != nil {
			log.Println("daemon", dro.Name, dro.Path, "recover", r)
		}
	}
	done := make(chan bool)
	outp := tools.ChangeExtension(dro.Path, "out.log")
	errp := tools.ChangeExtension(dro.Path, "err.log")
	envp := tools.ChangeExtension(dro.Path, "env")
	go func() {
		defer log.Println("daemon", dro.Name, dro.Path, "exited")
		defer traceRecover()
		defer close(done)
		run := func() {
			defer traceRecover()
			ff := os.O_APPEND | os.O_WRONLY | os.O_CREATE
			outf, err := os.OpenFile(outp, ff, 0644)
			tools.PanicIfError(err)
			defer outf.Close()
			errf, err := os.OpenFile(errp, ff, 0644)
			tools.PanicIfError(err)
			defer errf.Close()
			env := tools.ReadEnviron(envp)
			log.Println("daemon", dro.Name, dro.Name, dro.Path, "env", env)
			cmd := exec.Command(dro.Path)
			cmd.Env = env
			cmd.Stdout = outf
			cmd.Stderr = errf
			sin, err := cmd.StdinPipe()
			tools.PanicIfError(err)
			defer sin.Close()
			err = cmd.Start()
			tools.PanicIfError(err)
			pid := cmd.Process.Pid
			log.Println("daemon", dro.Name, dro.Path, "pid", pid)
			go func() {
				defer traceRecover()
				defer sin.Close()
				select {
				case <-exit:
				case <-done:
				}
			}()
			err = cmd.Wait()
			tools.PanicIfError(err)
		}
		count := 0
		for {
			if count > 0 {
				timer := time.NewTimer(2 * time.Second)
				select {
				case <-exit:
					timer.Stop()
					return
				case <-timer.C:
				}
			}
			log.Println("daemon", dro.Name, dro.Path, "run", count)
			run()
			count++
			select {
			case <-exit:
				return
			default:
				continue
			}
		}
	}()
	return done
}
