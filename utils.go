package main

import (
	"fmt"
	"math/rand"
)

func gen_random_user() string {
	alpha := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := string(alpha[rand.Intn(len(alpha)-10)])
	for i := 0; i < 9; i++ {
		result += string(alpha[rand.Intn(len(alpha))])
	}
	return result
}

func (server irc_server) get_formatted_server() string {
	return fmt.Sprintf("%s:%d", server.ip, server.port)
}

func ip_convert(ip int) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func unit_convert(bytes int) (int, string) {
	kb := bytes / 1000
	if kb < 1000 {
		return kb, "KB"
	}
	mb := bytes / 1000000
	if mb < 1000 {
		return mb, "MB"
	}
	gb := bytes / 1000000000
	return gb, "GB"
}
