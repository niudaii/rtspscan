package rtsp

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

type IPAddr struct {
	IP   string
	Port string
}

func (r *Runner) CheckAlive(targets []string) (results []IPAddr) {
	var tasks []IPAddr
	for _, target := range targets {
		parts := strings.Split(target, ":")
		if len(parts) != 2 {
			continue
		}
		tasks = append(tasks, IPAddr{
			IP:   parts[0],
			Port: parts[1],
		})
	}

	wg := &sync.WaitGroup{}
	rwMutex := &sync.RWMutex{}
	taskChan := make(chan IPAddr, r.options.Threads)

	for i := 0; i < r.options.Threads; i++ {
		go func() {
			for task := range taskChan {
				if r.conn(task) {
					rwMutex.Lock()
					results = append(results, task)
					rwMutex.Unlock()
				}
				wg.Done()
			}
		}()
	}

	bar := pb.StartNew(len(tasks))
	for _, task := range tasks {
		bar.Increment()
		wg.Add(1)
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()
	bar.Finish()

	return
}

func (r *Runner) conn(ipAddr IPAddr) bool {
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ipAddr.IP, ipAddr.Port), r.options.Timeout)
	return err == nil
}
