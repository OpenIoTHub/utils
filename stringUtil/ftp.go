package stringUtil

import (
	"errors"
	"strconv"
	"strings"
)

func DecodePassiveHostPort(line string) (host string, port, portPart1, portPart2 int, err error) {
	// PASV response format : 227 Entering Passive Mode (h1,h2,h3,h4,p1,p2).
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	if start == -1 || end == -1 {
		err = errors.New("invalid PASV response format")
		return
	}

	// We have to split the response string
	pasvData := strings.Split(line[start+1:end], ",")

	if len(pasvData) < 6 {
		err = errors.New("invalid PASV response format")
		return
	}

	// Let's compute the port number
	portPart1, err1 := strconv.Atoi(pasvData[4])
	if err1 != nil {
		err = err1
		return
	}

	portPart2, err2 := strconv.Atoi(pasvData[5])
	if err2 != nil {
		err = err2
		return
	}

	// Recompose port
	port = portPart1*256 + portPart2

	// Make the IP address to connect to
	host = strings.Join(pasvData[0:4], ".")
	return
}
