package main

import (
	"fmt"
	"log"
	"sync"
)

type managerDso struct {
	dao     Dao
	mutex   *sync.Mutex
	daemons map[string]*daemonDso
}

type daemonDso struct {
	dro  *DaemonDro
	exit chan bool
	done chan bool
}

type Manager interface {
	Close()
	Start(dro *DaemonDro) error
	Stop(name string) error
}

func NewManager(dao Dao) Manager {
	dso := &managerDso{}
	dso.dao = dao
	dso.mutex = &sync.Mutex{}
	dso.daemons = make(map[string]*daemonDso)
	list := dao.ListDaemons()
	for _, dro := range *list {
		log.Println("init", dro.Name, dro.Enabled, dro.Path)
		if dro.Enabled {
			err := dso.Start(&dro)
			if err != nil {
				dso.Close()
				log.Fatal(err)
			}
		}
	}
	return dso
}

func (dso *managerDso) Close() {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	for _, daemon := range dso.daemons {
		dro := daemon.dro
		log.Println("stop", dro.Name, dro.Enabled, dro.Path)
		close(daemon.exit)
		<-daemon.done
	}
	dso.daemons = make(map[string]*daemonDso)
}

func (dso *managerDso) Start(dro *DaemonDro) error {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	_, ok := dso.daemons[dro.Name]
	if ok {
		return fmt.Errorf("already started")
	}
	log.Println("start", dro.Name, dro.Enabled, dro.Path)
	daemon := &daemonDso{}
	daemon.dro = dro
	daemon.exit = make(chan bool)
	daemon.done = Run(dro, daemon.exit)
	dso.daemons[dro.Name] = daemon
	return nil
}

func (dso *managerDso) Stop(name string) error {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	daemon, ok := dso.daemons[name]
	if !ok {
		return fmt.Errorf("not running")
	}
	dro := daemon.dro
	log.Println("stop", dro.Name, dro.Enabled, dro.Path)
	delete(dso.daemons, name)
	close(daemon.exit)
	<-daemon.done
	return nil
}
