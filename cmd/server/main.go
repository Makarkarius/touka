package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"proj/internal/app/requester/comtradeapi"
	"proj/internal/app/server"
	"proj/internal/app/storage/mongostorage"
	"syscall"
)

const (
	defaultCfgPath     = "/etc/touka/"
	defaultCfgFilename = "production"
)

var (
	requesterCfgPath string
	storageCfgPath   string
	serverCfgPath    string
)

func init() {
	flag.StringVar(&requesterCfgPath,
		"requesterCfg",
		defaultCfgPath+"requester/"+defaultCfgFilename,
		"set path to requester config file")

	flag.StringVar(&storageCfgPath,
		"storageCfg",
		defaultCfgPath+"storage/"+defaultCfgFilename,
		"set path to storage config file")

	flag.StringVar(&serverCfgPath,
		"serverCfg",
		defaultCfgPath+"server/"+defaultCfgFilename,
		"set path to server config file")
}

func readCfg(target interface{}, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't open file\npath: %s\nerror: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read data\npath: %s\nerror: %w", path, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("can't unmarshal data\npath: %s\nerror: %w", path, err)
	}
	return nil
}

func main() {
	flag.Parse()
	var (
		serverCfg    server.Config
		requesterCfg comtradeapi.Config
		storageCfg   mongostorage.Config
	)
	if err := readCfg(&serverCfg, serverCfgPath); err != nil {
		fmt.Printf("can't read server config: %s\n", err)
		return
	}
	if err := readCfg(&requesterCfg, requesterCfgPath); err != nil {
		fmt.Printf("can't read requester config: %s\n", err)
		return
	}
	if err := readCfg(&storageCfg, storageCfgPath); err != nil {
		fmt.Printf("can't read storage config: %s\n", err)
		return
	}
	serverCfg.RequesterCfg = requesterCfg
	serverCfg.StorageCfg = storageCfg

	s, err := serverCfg.Build()
	if err != nil {
		fmt.Printf("can't build server from config: %s\n%v\n", err, serverCfg)
		return
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		<-sigint

		if err := s.Shutdown(context.Background()); err != nil {
			fmt.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	if err := s.Run(); err != nil {
		fmt.Println(err)
		return
	}
	<-idleConnectionsClosed
}
