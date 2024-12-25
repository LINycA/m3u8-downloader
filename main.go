package main

import (
	"m3u8-downloader/dot"
	"m3u8-downloader/utils"
	workprocess "m3u8-downloader/work_process"
)

func init() {
	preRun()
}

func main() {
	m3u8Uri := utils.Input("请输入完整的m3u8地址")
	filmName := utils.Input("请输入电影名称")
	workprocess.DownloadM3U8(m3u8Uri, filmName)
}

func preRun() {
	dot.SetupLogrus()
	dot.InitViper()
	utils.CheckFFmpegPath()
}
