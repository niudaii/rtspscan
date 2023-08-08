package rtsp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	StatusOK           = 200
	StatusUnauthorized = 401
	StatusNotFound     = 404
)

func DescribeMsg(url string, seq int, ua string) string {
	msgRet := "DESCRIBE " + url + " RTSP/1.0\r\n"
	msgRet += "CSeq: " + strconv.Itoa(seq) + "\r\n"
	if ua != "" {
		msgRet += "User-Agent: " + ua + "\r\n"
	}
	msgRet += "\r\n"
	return msgRet
}

func (r *Runner) Handler(serv Service) (status int, err error) {
	addr := fmt.Sprintf("%v:%v", serv.IP, serv.Port)
	conn, err := net.DialTimeout("tcp", addr, r.options.Timeout)
	if err != nil {
		return
	}
	_ = conn.SetReadDeadline(time.Now().Add(r.options.Timeout))
	_ = conn.SetWriteDeadline(time.Now().Add(r.options.Timeout))

	// just omit OPTIONS, send DESCRIBE
	seq := 1
	data := make([]byte, 255)
	seq += 1
	msg := DescribeMsg(serv.URL, seq, r.options.UserAgent)
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return
	}
	_, err = conn.Read(data)
	if err != nil {
		return
	}
	if strings.Contains(string(data), "200 OK") {
		status = StatusOK
		return
	}
	if strings.Contains(string(data), "401") {
		status = StatusUnauthorized
		return
	}
	if strings.Contains(string(data), "404") {
		status = StatusNotFound
		return
	}
	return
}
