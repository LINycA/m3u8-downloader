package workprocess

import (
	"m3u8-downloader/utils"
	progressbar "m3u8-downloader/utils/progress_bar"
	"m3u8-downloader/utils/request"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/semaphore"
)

// TsFileDownloader 下载ts媒体文件
func TsFileDownloader(uri, filePath string, cli *http.Client, headers map[string]string, wg *sync.WaitGroup, sem *semaphore.Weighted, bar *progressbar.ProgressBar) {
	defer wg.Done()
	defer sem.Release(1)
	if utils.CheckTsDirExists(filePath) {
		bar.Add(1)
		return
	} else {
		tsFile, _ := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		defer tsFile.Close()
		req := request.NewRequest(http.MethodGet, "", uri, 5, headers, cli)
		req.DoRequest()
		tsFile.Write(req.ResponseBody)
		bar.Add(1)
	}
}
