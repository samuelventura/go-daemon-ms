package main

import (
	"log"
	"os"

	"github.com/samuelventura/go-state"
	"github.com/samuelventura/go-tools"
	"github.com/samuelventura/go-tree"
)

func entry(inter bool, exit chan bool) {
	ctrlc := tools.SetupCtrlc()
	stdin := tools.SetupStdinAll()

	log.Println("start", os.Getpid())
	defer log.Println("exit")

	rnode := tree.NewRoot("root", log.Println)
	defer rnode.WaitDisposed()
	//recover closes as well
	defer rnode.Recover()

	spath := tools.WithExtension("state")
	snode := state.Serve(rnode, spath)
	defer snode.WaitDisposed()
	defer snode.Close()
	log.Println("socket", spath)

	anode := rnode.AddChild("api")
	defer anode.WaitDisposed()
	defer anode.Close()
	tools.LoadDefaultEnviron()
	anode.SetValue("driver", tools.GetEnviron("DAEMON_DB_DRIVER", "sqlite"))
	anode.SetValue("source", tools.GetEnviron("DAEMON_DB_SOURCE", tools.WithExtension("db3")))
	anode.SetValue("endpoint", tools.GetEnviron("DAEMON_ENDPOINT", "127.0.0.1:31600"))
	dao := NewDao(anode) //close on root
	rnode.AddCloser("dao", dao.Close)
	anode.SetValue("dao", dao)
	api(anode)

	select {
	case <-rnode.Closed():
	case <-snode.Closed():
	case <-anode.Closed():
	case <-ctrlc:
	case <-stdin:
	case <-exit:
	}
}
