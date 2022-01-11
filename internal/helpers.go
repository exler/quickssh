package internal

import (
	"bufio"
	"log"
	"os"
	"strings"

	"golang.org/x/term"
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

func GetTerminalSize() (width, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0
	}
	return
}
