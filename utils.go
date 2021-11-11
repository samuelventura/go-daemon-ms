package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadenv() {
	envp := withext(".env")
	env := environ(envp)
	for _, line := range env {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			err := fmt.Errorf("invalid env %s", line)
			log.Fatal(err)
		}
		os.Setenv(parts[0], parts[1])
	}
}

func getenv(name string, defval string) string {
	value := os.Getenv(name)
	trimmed := strings.TrimSpace(value)
	if len(trimmed) > 0 {
		log.Println(name, value)
		return value
	}
	log.Println(name, defval)
	return defval
}

func withext(ext string) string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	base := filepath.Base(exe)
	file := base + "." + ext
	return filepath.Join(dir, file)
}

func changeext(path string, next string) string {
	ext := filepath.Ext(path) //includes .
	npath := strings.TrimSuffix(path, ext)
	return npath + next
}

func environ(path string) []string {
	lines := []string{}
	file, err := os.Open(path)
	if err != nil {
		return lines
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			lines = append(lines, trimmed)
		}
	}
	return lines
}
