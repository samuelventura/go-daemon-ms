package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/samuelventura/go-tree"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type daoDso struct {
	mutex *sync.Mutex
	db    *gorm.DB
}

type Dao interface {
	Close() error
	ListDaemons() []*DaemonDro
	GetDaemon(name string) (*DaemonDro, error)
	CreateDaemon(name string, path string) (*DaemonDro, error)
	EnableDaemon(name string, enabled bool) error
	DelDaemon(name string) error
}

func Dialector(node tree.Node) gorm.Dialector {
	driver := node.GetValue("driver").(string)
	source := node.GetValue("source").(string)
	switch driver {
	case "sqlite":
		return sqlite.Open(source)
	case "postgres":
		return postgres.Open(source)
	}
	log.Fatalf("unknown driver %s", driver)
	return nil
}

func NewDao(node tree.Node) Dao {
	mode := logger.Default.LogMode(logger.Silent)
	config := &gorm.Config{Logger: mode}
	dialector := Dialector(node)
	db, err := gorm.Open(dialector, config)
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&DaemonDro{})
	if err != nil {
		log.Fatal(err)
	}
	return &daoDso{&sync.Mutex{}, db}
}

func (dso *daoDso) Close() error {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	db, err := dso.db.DB()
	if err != nil {
		log.Fatal(err)
	}
	return db.Close()
}

func (dso *daoDso) ListDaemons() []*DaemonDro {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	dros := []*DaemonDro{}
	result := dso.db.Where("true").Find(&dros)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return dros
}

func (dso *daoDso) GetDaemon(name string) (*DaemonDro, error) {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	dro := &DaemonDro{}
	result := dso.db.Where("name = ?", name).First(dro)
	return dro, result.Error
}

func (dso *daoDso) CreateDaemon(name string, path string) (*DaemonDro, error) {
	if len(strings.TrimSpace(path)) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	count := int64(0)
	dso.db.Model(&DaemonDro{}).
		Where("name = ?", name).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("duplicated name")
	}
	dro := &DaemonDro{}
	dro.Name = name
	dro.Path = path
	result := dso.db.Create(dro)
	return dro, result.Error
}

func (dso *daoDso) EnableDaemon(name string, enabled bool) error {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	result := dso.db.Model(&DaemonDro{}).
		Where("name = ?", name).Update("Enabled", enabled)
	if result.Error == nil && result.RowsAffected != 1 {
		return fmt.Errorf("daemon not found")
	}
	return result.Error
}

func (dso *daoDso) DelDaemon(name string) error {
	dso.mutex.Lock()
	defer dso.mutex.Unlock()
	result := dso.db.Where("name = ?", name).
		Delete(&DaemonDro{})
	if result.Error == nil && result.RowsAffected != 1 {
		return fmt.Errorf("daemon not found")
	}
	return result.Error
}
