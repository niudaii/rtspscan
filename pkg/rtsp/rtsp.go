package rtsp

import (
	"net"
	"strconv"
	"strings"
)

const (
	StatusOK           = 200
	StatusUnauthorized = 401
	StatusNotFound     = 404
)

func (r *Runner) Handler(conn net.Conn, serv Service) (status int, err error) {
	// just omit OPTIONS, send DESCRIBE
	seq := 1
	data := make([]byte, 255)
	seq += 1
	msg := describeMsg(serv.URL, seq, r.options.UserAgent)
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

func describeMsg(url string, seq int, ua string) string {
	msgRet := "DESCRIBE " + url + " RTSP/1.0\r\n"
	msgRet += "CSeq: " + strconv.Itoa(seq) + "\r\n"
	if ua != "" {
		msgRet += "User-Agent: " + ua + "\r\n"
	}
	msgRet += "\r\n"
	return msgRet
}
