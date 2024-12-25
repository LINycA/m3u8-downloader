package request_test

import (
	"fmt"
	"m3u8-downloader/utils/request"
	"net/http"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	cli, err := request.NewClient(time.Second*10, "socks5://127.0.0.1:7890")
	if err != nil {
		fmt.Println(err)
		return
	}
	req := request.NewRequest(http.MethodGet, "", "https://www.baidu.com", 1, nil, cli)
	err = req.DoRequest()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(req.ResponseBody))
}
