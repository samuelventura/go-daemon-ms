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

	log.Println("start", os.Getpid())
	defer log.Println("exit")

	rnode := tree.NewRoot("root", log.Println)
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
	loadenv()
	anode.SetValue("driver", getenv("DAEMON_DB_DRIVER", "sqlite"))
	anode.SetValue("source", getenv("DAEMON_DB_SOURCE", withext("db3")))
	anode.SetValue("endpoint", getenv("DAEMON_ENDPOINT", "127.0.0.1:31600"))
	dao := NewDao(anode) //close on root
	rnode.AddCloser("dao", dao.Close)
	anode.SetValue("dao", dao)
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
