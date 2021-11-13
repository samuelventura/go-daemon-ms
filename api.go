package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samuelventura/go-tree"
)

func api(node tree.Node) {
	dao := node.GetValue("dao").(Dao)
	manager := NewManager(dao)
	node.AddAction("manager", manager.Close)
	endpoint := node.GetValue("endpoint").(string)
	gin.SetMode(gin.ReleaseMode) //remove debug warning
	router := gin.New()          //remove default logger
	router.Use(gin.Recovery())   //looks important
	rapi := router.Group("/api/daemon")
	rapi.GET("/list", func(c *gin.Context) {
		list := dao.ListDaemons()
		c.JSON(200, list)
	})
	rapi.GET("/info/:name", func(c *gin.Context) {
		name := c.Param("name")
		row, err := dao.GetDaemon(name)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, row)
	})
	rapi.GET("/env/:name", func(c *gin.Context) {
		name := c.Param("name")
		row, err := dao.GetDaemon(name)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		envp := changeext(row.Path, ".env")
		c.JSON(200, map[string]interface{}{"Path": envp, "Vars": environ(envp)})
	})
	rapi.DELETE("/env/:name", func(c *gin.Context) {
		name := c.Param("name")
		row, err := dao.GetDaemon(name)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		envp := changeext(row.Path, ".env")
		err = os.Remove(envp)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, row)
	})
	rapi.POST("/env/:name", func(c *gin.Context) {
		name := c.Param("name")
		row, err := dao.GetDaemon(name)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		envp := changeext(row.Path, ".env")
		ff := os.O_TRUNC | os.O_WRONLY | os.O_CREATE
		envf, err := os.OpenFile(envp, ff, 0644)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		defer envf.Close()
		key := http.CanonicalHeaderKey("DaemonEnviron")
		values := c.Request.Header[key]
		for _, line := range values {
			fmt.Fprintln(envf, strings.TrimSpace(line))
		}
		c.JSON(200, gin.H{
			"Path": envp,
			"Vars": environ(envp)})
	})
	rapi.POST("/install/:name", func(c *gin.Context) {
		name := c.Param("name")
		path, _ := c.GetQuery("path")
		row, err := dao.CreateDaemon(name, path)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, row)
	})
	rapi.POST("/enable/:name", func(c *gin.Context) {
		name := c.Param("name")
		err := dao.EnableDaemon(name, true)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, "ok")
	})
	rapi.POST("/disable/:name", func(c *gin.Context) {
		name := c.Param("name")
		err := dao.EnableDaemon(name, false)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, "ok")
	})
	//does not stop the daemon
	rapi.POST("/uninstall/:name", func(c *gin.Context) {
		name := c.Param("name")
		err := dao.DelDaemon(name)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, "ok")
	})
	rapi.POST("/start/:name", func(c *gin.Context) {
		name := c.Param("name")
		row, err := dao.GetDaemon(name)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		err = manager.Start(row)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, "ok")
	})
	//does not check if daemon exists
	rapi.POST("/stop/:name", func(c *gin.Context) {
		name := c.Param("name")
		err := manager.Stop(name)
		if err != nil {
			c.JSON(400, fmt.Sprintf("err: %v", err))
			return
		}
		c.JSON(200, "ok")
	})
	listen, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatal(err)
	}
	node.AddCloser("listen", listen.Close)
	port := listen.Addr().(*net.TCPAddr).Port
	log.Println("port", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: router,
	}
	node.AddProcess("server", func() {
		err = server.Serve(listen)
		if err != nil {
			log.Println(endpoint, port, err)
		}
	})
}
