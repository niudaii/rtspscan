package rtsp

import (
	"fmt"
	"github.com/zp857/goutil/networkx"
	"net"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/niudaii/util"
)

type Options struct {
	UserAgent string
	Threads   int
	PathList  []string
	UserList  []string
	PassList  []string
	Proxy     string
	Timeout   time.Duration
}

type Runner struct {
	options *Options
}

func NewRunner(options *Options) (*Runner, error) {
	return &Runner{
		options: options,
	}, nil
}

type Service struct {
	IP   string
	Port string
	Path string
	URL  string
}

func (r *Runner) Run(targets []IPAddr) (results []string) {
	var tasks []Service
	for _, target := range targets {
		tasks = append(tasks, Service{
			IP:   target.IP,
			Port: target.Port,
		})
	}

	wg := &sync.WaitGroup{}
	rwMutex := &sync.RWMutex{}
	taskChan := make(chan Service, r.options.Threads)

	for i := 0; i < r.options.Threads; i++ {
		go func() {
			for task := range taskChan {
				result, err := r.Scan(task)
				if err != nil {
					log.Debug(err.Error())
				} else {
					if result.URL != "" {
						log.Printf("[+] %v", result.URL)
						rwMutex.Lock()
						results = append(results, result.URL)
						rwMutex.Unlock()
					}
				}
				wg.Done()
			}
		}()
	}

	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()

	return
}

func (r *Runner) Scan(serv Service) (result Service, err error) {
	addr := fmt.Sprintf("%v:%v", serv.IP, serv.Port)
	var conn net.Conn
	conn, err = networkx.NewConn(addr, r.options.Proxy, r.options.Timeout)
	if err != nil {
		return
	}
	r.options.PathList = append([]string{""}, r.options.PathList...)
	r.options.PathList = util.RemoveDuplicate(r.options.PathList)
	// check path
	var status int
	for _, path := range r.options.PathList {
		serv.Path = path
		serv.URL = fmt.Sprintf("rtsp://%v:%v%v", serv.IP, serv.Port, path)
		status, err = r.Handler(conn, serv)
		if err != nil {
			return
		}
		if status == StatusOK {
			result = serv
			return
		}
		if status == StatusNotFound {
			continue
		}
		if status == StatusUnauthorized {
			break
		}
	}
	if status == StatusNotFound {
		err = fmt.Errorf("path not found")
		return
	}
	// brute auth
	for _, user := range r.options.UserList {
		for _, pass := range r.options.PassList {
			serv.URL = fmt.Sprintf("rtsp://%v:%v@%v:%v%v", user, pass, serv.IP, serv.Port, serv.Path)
			status, err = r.Handler(conn, serv)
			if err != nil {
				return
			}
			if status == StatusOK {
				result = serv
				return
			}
			if status == StatusUnauthorized {
				continue
			}
		}
	}
	if status == StatusUnauthorized {
		err = fmt.Errorf("auth not found")
	}
	return
}
