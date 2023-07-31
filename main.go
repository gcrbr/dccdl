package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Print("DCC-DL ;; created by @gcrbr\n\n")

	server := irc_server{"0.0.0.0", 6667}
	var channel string
	var bot string
	file_id := -1
	save := false

	flag.StringVar(&server.ip, "i", "", "IRC server IP address")
	flag.IntVar(&server.port, "p", 6667, "IRC server port")
	flag.StringVar(&channel, "c", "", "Channel name")
	flag.StringVar(&bot, "b", "", "Bot name")
	flag.IntVar(&file_id, "f", 0, "File ID")
	flag.BoolVar(&save, "s", false, "Save to file")
	flag.Parse()

	if server.ip == "" || channel == "" || bot == "" || file_id == -1 {
		fmt.Println("Missing arguments")
		return
	}

	file := server.xdcc_download(channel, bot, file_id)

	if file == (dcc_connection{}) {
		return
	}

	size, unit := unit_convert(file.filesize)

	fmt.Printf("\nFile name	%s\n", file.filename)
	fmt.Printf("Address		%s:%d\n", file.ip, file.port)
	fmt.Printf("Size		%d %s\n\n", size, unit)

	if save {
		percent := 0.0
		written := 0
		f, _ := os.OpenFile(file.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		con, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", file.ip, file.port))
		tmp := make([]byte, 1)
		var err error = nil
		for err == nil {
			_, err = con.Read(tmp)
			f.Write(tmp)
			written += 1
			percent = (float64(written) / float64(file.filesize)) * 100
			if percent > 100 {
				percent = 100
			}
			fmt.Printf("\rDownloaded %.2f%%", percent)
		}
		if err != nil && err != io.EOF {
			fmt.Printf("\nDownload error: %s", err)
		}
		f.Close()
	}

}
