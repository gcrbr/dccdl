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
	start := time.Now()
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
			if time.Since(start).Seconds() >= 60 {
				fmt.Println("Timeout")
				return dcc_connection{}
			}
			if strings.HasPrefix(line, "PING :") {
				socket.Write([]byte(strings.Replace(line, "PING", "PONG", 1)))
			}

			words := strings.Split(line, " ")

			//fmt.Printf("xxxxxx->%s\n", line)

			if lock {
				if len(words) > 0 && words[1] == "001" {
					lock = false
					fmt.Println("Session established")
					socket.Write([]byte(fmt.Sprintf("JOIN #%s\r\n", channel)))
					socket.Write([]byte(fmt.Sprintf("PRIVMSG %s :xdcc send #%d\r\n", bot, file_id)))
					start = time.Now()
					fmt.Println("Requesting file...")
				}
			} else {
				if len(words) > 3 && words[1] == "PRIVMSG" && strings.EqualFold(words[2], username) && words[3] == ":\x01DCC" && words[4] == "SEND" {
					socket.Write([]byte(fmt.Sprintf("PRIVMSG %s :DCC ACCEPT %s %s 0\r\n", bot, words[5], words[7])))
					output.filename = words[5]
					ip, err2 := strconv.Atoi(words[6])
					if err2 != nil {
						fmt.Printf("Invalid IP address '%s'\n", words[6])
						return dcc_connection{}
					}
					output.ip = ip_convert(ip)
					output.port, err2 = strconv.Atoi(words[7])
					if err2 != nil {
						fmt.Printf("Invalid port number '%s'\n", words[7])
						return dcc_connection{}
					}
					output.filesize, err2 = strconv.Atoi(strings.Split(words[8], "\x01")[0])
					if err2 != nil {
						fmt.Printf("Invalid file size '%s'\n", words[8])
						return dcc_connection{}
					}
					break
				}
				if words[1] == "NOTICE" && strings.EqualFold(words[2], username) {
					fmt.Printf("Message from server: %s\n", strings.Join(words[3:], " ")[1:])
				}
			}
			line = ""
		}

	}
	return output
}
