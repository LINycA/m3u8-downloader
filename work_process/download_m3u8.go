package workprocess

import (
	"bytes"
	"context"
	"fmt"
	"m3u8-downloader/dot"
	"m3u8-downloader/utils"
	progressbar "m3u8-downloader/utils/progress_bar"
	"m3u8-downloader/utils/request"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
)

// DownloadM3U8 下载m3u8文件
func DownloadM3U8(uri, filmName string) {
	m3u8Url, err := url.Parse(uri)
	if err != nil {
		dot.Logger().WithError(err).WithField("url", uri).Error("m3u8地址解析")
		return
	}
	headers := map[string]string{
		"cookie":     viper.GetString("cookie"),
		"user-agent": viper.GetString("ua"),
	}
	cli, err := request.NewClient(time.Minute*5, viper.GetString("proxy"))
	if err != nil {
		dot.Logger().WithError(err).Error("创建http客户端")
		return
	}
	req := request.NewRequest(http.MethodGet, "", m3u8Url.String(), 3, headers, cli)
	err = req.DoRequest()
	if err != nil {
		dot.Logger().WithError(err).Error("请求时错误")
		return
	}

	timeout := time.Minute * 5
	tsCli := &http.Client{
		Timeout: timeout,
	}

	m3u8FileBuffer := bytes.NewBuffer(req.ResponseBody)
	m3u8Pl, m3u8Type, err := m3u8.Decode(*m3u8FileBuffer, true)
	if err != nil {
		dot.Logger().WithError(err).Error("m3u8文件解析")
		return
	}
	wkdir, err := os.Getwd()
	if err != nil {
		dot.Logger().WithError(err).Panic("检测文件,获取当前工作文件夹")
		return
	}
	tsDir := path.Join(wkdir, filmName, "ts_file")
	m3u8FileContent := ""
	utils.CheckDir(tsDir)
	switch m3u8Type {
	case m3u8.MASTER:
		// 这个m3u8是主文件，底下还有其他子地址,一般是按照清晰度命名，如有其他形式再优化
		playlist := m3u8Pl.(*m3u8.MasterPlaylist)
		variantsMap := map[int]*m3u8.Variant{}
		variantsKeyL := []int{}
		for _, v := range playlist.Variants {
			ind, err := strconv.Atoi(v.Name)
			if err != nil {
				dot.Logger().WithError(err).Error("m3u8主文件解析错误")
				return
			}
			variantsMap[ind] = v
			variantsKeyL = append(variantsKeyL, ind)
		}
		utils.IntListReverse(variantsKeyL)
		DownloadM3U8(variantsMap[variantsKeyL[0]].URI, filmName)
	case m3u8.MEDIA:
		playlist := m3u8Pl.(*m3u8.MediaPlaylist)
		total := 0
		for _, seg := range playlist.Segments {
			if seg != nil {
				total += 1
			}
		}
		reg, _ := regexp.Compile(`[0-9a-zA-Z_-]+?\.ts.*`)
		tsList := reg.FindAllString(string(req.ResponseBody), -1)
		bar := progressbar.NewProgressBar(filmName, len(tsList))
		wg := &sync.WaitGroup{}
		sem := semaphore.NewWeighted(viper.GetInt64("semaphore"))
		pathLastNum := strings.LastIndex(m3u8Url.Path, "/")
		for ind, seg := range playlist.Segments {
			if seg != nil {
				tsFilePath := path.Join(tsDir, fmt.Sprintf("%d.ts", ind))
				m3u8FileContent += fmt.Sprintf("file '%s'\n", tsFilePath)
				wg.Add(1)
				sem.Acquire(context.Background(), 1)
				tsUri := fmt.Sprintf("%s://%s%s%s",
					m3u8Url.Scheme,
					m3u8Url.Host,
					m3u8Url.Path[:pathLastNum+1],
					tsList[ind],
				)
				go TsFileDownloader(tsUri, tsFilePath, tsCli, headers, wg, sem, bar)
			}
		}
		wg.Wait()
		m3u8FilePath := path.Join(wkdir, filmName, "merge.m3u8")
		m3u8File, err := os.OpenFile(m3u8FilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			dot.Logger().WithError(err).Error("保存m3u8文件")
			return
		}
		_, err = m3u8File.WriteString(m3u8FileContent)
		if err != nil {
			dot.Logger().WithError(err).Error("m3u8文件写入错误")
			return
		}
		m3u8File.Close()
		time.Sleep(time.Second)
		outputFilePath := path.Join(wkdir, fmt.Sprintf("%s.mp4", filmName))
		utils.FFmpegMergeTs(m3u8FilePath, outputFilePath)
		dot.Logger().WithField("影片", filmName).Info("影片下载完成")
	}
}
