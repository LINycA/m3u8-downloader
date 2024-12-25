package utils

import (
	"bytes"
	"fmt"
	"m3u8-downloader/dot"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// CheckTsDirExists 检测ts文件是否存在是否下载完整
func CheckTsDirExists(tsFilePath string) bool {
	if fileInfo, err := os.Stat(tsFilePath); os.IsNotExist(err) || fileInfo.Size() < 10 {
		return false
	}
	return true
}

// CheckDir 检测文件夹是否存在，不存在则创建
func CheckDir(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}
}

// CheckFFmpegPath 检测系统中是否有ffmpeg工具
func CheckFFmpegPath() {
	runTimeOs := runtime.GOOS
	findTool := "which"
	switch runTimeOs {
	case "darwin":
		findTool = "which"
	case "linux":
		findTool = "which"
	case "windows":
		findTool = "where"
	}
	replacer := strings.NewReplacer("\n", "", "|", "", "\t", "", "\r", "")
	var stdout bytes.Buffer
	cmd := exec.Command(findTool, "ffmpeg")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		dot.Logger().WithError(err).Panic("检测ffmpeg环境")
	}
	ffmpegPath := replacer.Replace(stdout.String())
	viper.SetDefault("ffmpeg_path", ffmpegPath)
}

// FFmpegMergeTs 通过m3u8文件记录的ts文件顺序与位置进行ts文件拼接
func FFmpegMergeTs(m3u8FilePath string, outputPath string) {
	ffmpegPath := viper.GetString("ffmpeg_path")
	fmt.Println(ffmpegPath)
	cmd := exec.Command(ffmpegPath, "-f", "concat", "-safe", "0", "-i", m3u8FilePath, "-c", "copy", outputPath)
	err := cmd.Run()
	if err != nil {
		dot.Logger().WithError(err).Error("ts文件合并")
	}
}
