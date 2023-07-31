package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type irc_server struct {
	ip   string
	port int
}

type dcc_connection struct {
	filename string
	ip       string
	port     int
	filesize int
}

var username string

func (server irc_server) xdcc_download(channel string, bot string, file_id int) dcc_connection {
	var output dcc_connection
	var start time.Time
	lock := true

	fmt.Println("Connecting...")
	socket, err1 := net.Dial("tcp", server.get_formatted_server())

	if err1 != nil {
		fmt.Printf("Connection error: %s\n", err1)
		return dcc_connection{}
	}

	fmt.Println("Connected")
	username = gen_random_user()
	socket.Write([]byte(fmt.Sprintf("NICK %s\r\n", username)))
	socket.Write([]byte(fmt.Sprintf("USER %s 8 * : %s\r\n", username, username)))

	line := ""
	for {
		tmp := make([]byte, 1)
		socket.Read(tmp)
		s_tmp := string(tmp)
		if s_tmp != "\n" {
			line += s_tmp
		} else {
			if time.Since(start).Seconds() >= 90 { //1 min
				fmt.Println("Timeout")
				return dcc_connection{}
			}
			if strings.HasPrefix(line, "PING :") {
				socket.Write([]byte(strings.Replace(line, "PING", "PONG", 1)))
			}
			if lock {
				if strings.Split(line, " ")[1] == "266" {
					lock = false
					fmt.Println("Session established")
					socket.Write([]byte(fmt.Sprintf("JOIN #%s\r\n", channel)))
					socket.Write([]byte(fmt.Sprintf("PRIVMSG %s :xdcc send #%d\r\n", bot, file_id)))
					start = time.Now()
					fmt.Println("Requesting file...")
				}
			} else {
				if strings.Contains(line, "PRIVMSG") && strings.Contains(line, ":\x01DCC SEND") {
					dcc_data := strings.Split(line, " ")
					socket.Write([]byte(fmt.Sprintf("PRIVMSG %s :DCC ACCEPT %s %s 0\r\n", bot, dcc_data[5], dcc_data[7])))
					output.filename = dcc_data[5]
					ip, err2 := strconv.Atoi(dcc_data[6])
					if err2 != nil {
						fmt.Printf("Invalid IP address '%s'\n", dcc_data[6])
						return dcc_connection{}
					}
					output.ip = ip_convert(ip)
					output.port, err2 = strconv.Atoi(dcc_data[7])
					if err2 != nil {
						fmt.Printf("Invalid port number '%s'\n", dcc_data[7])
						return dcc_connection{}
					}
					output.filesize, err2 = strconv.Atoi(strings.Split(dcc_data[8], "\x01")[0])
					if err2 != nil {
						fmt.Printf("Invalid file size '%s'\n", dcc_data[8])
						return dcc_connection{}
					}
					break
				}
				if strings.Contains(line, username) && strings.Contains(line, bot) && strings.Split(line, " ")[1] != "353" && !strings.Contains(line, ":\x01DCC SEND") {
					if strings.Contains(line, "NOTICE") {
						fmt.Printf("Message from server: %s\n", strings.Split(line, "NOTICE "+username+" :** ")[1])
					} else {
						fmt.Println(line)
					}
				}
			}
			line = ""
		}

	}
	return output
}
