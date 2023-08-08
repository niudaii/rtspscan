package cmd

import (
	"flag"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/niudaii/util"
	"github.com/niudaii/util/files"
)

type Options struct {
	Input      string
	InputFile  string
	OutputFile string
	UserAgent  string
	PassFile   string
	UserFile   string
	PathFile   string
	Threads    int
	Timeout    int
	Debug      bool

	InputList []string
	PathList  []string
	UserList  []string
	PassList  []string
}

func InitOptions(options *Options) (err error) {
	flag.StringVar(&options.Input, "i", "", "The input")
	flag.StringVar(&options.InputFile, "f", "", "The input file")
	flag.StringVar(&options.OutputFile, "o", "output.txt", "The output file")
	flag.StringVar(&options.UserAgent, "ua", "", "The user agent")
	flag.StringVar(&options.PathFile, "path-file", "resource/path.txt", "The path file")
	flag.StringVar(&options.UserFile, "user-file", "resource/user.txt", "The user file")
	flag.StringVar(&options.PassFile, "pass-file", "resource/pass.txt", "The pass file")
	flag.IntVar(&options.Threads, "t", 100, "Thread num")
	flag.IntVar(&options.Timeout, "timeout", 5, "Timeout in seconds")
	flag.BoolVar(&options.Debug, "debug", false, "Enable debug mode")
	flag.Parse()

	if err = options.validateOptions(); err != nil {
		return
	}
	if err = options.configureOptions(); err != nil {
		return
	}
	options.configureOutput()
	options.printOptions()

	return
}

func (o *Options) validateOptions() (err error) {
	if o.Input == "" && o.InputFile == "" {
		return fmt.Errorf("no input provided")
	}
	if o.InputFile != "" && !files.FileExists(o.InputFile) {
		return fmt.Errorf("file %v does not exist", o.InputFile)
	}

	return
}

func (o *Options) configureOutput() {
	if o.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetReportTimestamp(false)
}

func (o *Options) configureOptions() (err error) {
	if o.Input != "" {
		o.InputList = append(o.InputList, o.Input)
	} else {
		if o.InputList, err = files.ReadLines(o.InputFile); err != nil {
			return
		}
	}
	if o.PathList, err = files.ReadLines(o.PathFile); err != nil {
		return
	}
	if o.UserList, err = files.ReadLines(o.UserFile); err != nil {
		return
	}
	if o.PassList, err = files.ReadLines(o.PassFile); err != nil {
		return
	}

	o.InputList = util.RemoveDuplicate(o.InputList)
	o.PathList = util.RemoveDuplicate(o.PathList)
	o.UserList = util.RemoveDuplicate(o.UserList)
	o.PassList = util.RemoveDuplicate(o.PassList)

	return nil
}

func (o *Options) printOptions() {
	log.Debug("The input list", "total", len(o.InputList))
	log.Debug("The path list", "total", len(o.PathList))
	log.Debug("The user list", "total", len(o.UserList))
	log.Debug("The pass list", "total", len(o.PassList))
}
