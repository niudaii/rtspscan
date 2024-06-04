package cmd

import (
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/niudaii/rtspscan/pkg/rtsp"
	"github.com/niudaii/util/files"
)

func Run() {
	var options Options
	err := initOptions(&options)
	if err != nil {
		log.Error(err.Error())
		return
	}
	rtspOptions := &rtsp.Options{
		UserAgent: options.UserAgent,
		Threads:   options.Threads,
		Timeout:   time.Duration(options.Timeout) * time.Second,
		Proxy:     options.Proxy,
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
	log.Info("Scan finished", "results", len(results), "spentTime", time.Since(start))
	err = files.WriteFile(options.OutputFile, strings.Join(results, "\n"))
	if err != nil {
		log.Error(err.Error())
	}
}
