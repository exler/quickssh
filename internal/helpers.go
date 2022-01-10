package internal

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func GetUserInput() (text string) {
	text, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	text = strings.Trim(text, "\r\n")
	return
}

func ParseUserAndHost(target string) (user, host string) {
	res := strings.Split(target, "@")
	if len(res) != 2 {
		log.Fatalln("Wrong format of the target!")
	}

	return res[0], res[1]
}
