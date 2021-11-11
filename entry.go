package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/samuelventura/go-state"
	"github.com/samuelventura/go-tree"
)

func entry(inter bool, exit chan bool) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	log.Println("start")
	defer log.Println("exit")

	loadenv()
	args := NewArgs()
	args.Set("driver", getenv("DAEMON_DB_DRIVER", "sqlite"))
	args.Set("source", getenv("DAEMON_DB_SOURCE", withext("db3")))
	args.Set("endpoint", getenv("DAEMON_ENDPOINT", "127.0.0.1:31600"))

	dao := NewDao(args)
	defer dao.Close()

	rlog := tree.NewLog()
	rnode := tree.NewRoot("root", rlog)
	defer rnode.WaitDisposed()
	//recover closes as well
	defer rnode.Recover()

	spath := state.SingletonPath()
	snode := state.Serve(rnode, spath)
	defer snode.WaitDisposed()
	defer snode.Close()
	log.Println("socket", spath)

	anode := rnode.AddChild("api")
	defer anode.WaitDisposed()
	defer anode.Close()
	anode.SetValue("dao", dao)
	anode.SetValue("endpoint", args.Get("endpoint"))
	api(anode)

	stdin := make(chan interface{})
	go func() {
		if inter {
			defer close(stdin)
			ioutil.ReadAll(os.Stdin)
		}
	}()
	select {
	case <-rnode.Closed():
	case <-snode.Closed():
	case <-anode.Closed():
	case <-ctrlc:
	case <-stdin:
	case <-exit:
	}
}
