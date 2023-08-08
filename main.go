package main

import (
	"rtspscan/internal/cmd"
	"rtspscan/pkg/rtsp"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/niudaii/util/files"
)

func main() {
	var options cmd.Options
	err := cmd.InitOptions(&options)
	if err != nil {
		log.Error(err.Error())
		return
	}
	rtspOptions := &rtsp.Options{
		UserAgent: options.UserAgent,
		Threads:   options.Threads,
		Timeout:   time.Duration(options.Timeout) * time.Second,
		PathList:  options.PathList,
		UserList:  options.UserList,
		PassList:  options.PassList,
	}
	runner, err := rtsp.NewRunner(rtspOptions)
	if err != nil {
		log.Error(err.Error())
		return
	}
	start := time.Now()
	log.Info("Start checking addr alive")
	aliveAddr := runner.CheckAlive(options.InputList)
	log.Info("Start scanning rtsp")
	results := runner.Run(aliveAddr)
	log.Info("Scan finished", "results", len(results), "spent", time.Since(start))
	err = files.WriteFile(options.OutputFile, strings.Join(results, "\n"))
	if err != nil {
		log.Error(err.Error())
		return
	}
}
