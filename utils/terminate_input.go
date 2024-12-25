package utils

import (
	"bufio"
	"fmt"
	"m3u8-downloader/dot"
	"os"
	"strings"
)

// Input 接收终端输入
func Input(tip string) string {
	fmt.Println(tip)
	stdin := bufio.NewReader(os.Stdin)
	stdinReplacer := strings.NewReplacer("\n", "", "|", "", "\r", "", "\t", "")
	text, err := stdin.ReadString('\n')
	if err != nil {
		dot.Logger().WithError(err).Panic("终端输入")
	}
	return stdinReplacer.Replace(text)
}
