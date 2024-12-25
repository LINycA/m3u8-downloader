package progressbar

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type ProgressBar struct {
	sync.Mutex
	Title        string   // 进度条标题
	Total        int      // 进度条总量
	Width        int      // 进度条在终端显示的长度
	Progress     int      // 进度条进度
	ProgressChan chan int // 进度条通道
}

// NewProgressBar 新建一个进度条
func NewProgressBar(title string, total int, width ...int) *ProgressBar {
	progressChan := make(chan int)
	var barWidth int
	if len(width) > 0 && width[0] >= 50 {
		barWidth = 50
	} else {
		barWidth = 50
	}
	progressBar := &ProgressBar{
		Total:        total,
		Title:        title,
		Width:        barWidth,
		ProgressChan: progressChan,
	}
	progressBar.start()
	time.Sleep(time.Millisecond * 100)
	progressBar.ProgressChan <- 0
	return progressBar
}

// start 进度条打印
func (p *ProgressBar) start() {
	go func() {
		ticker := time.NewTicker(time.Minute * 30)
		defer ticker.Stop()
		for {
			select {
			case progress := <-p.ProgressChan:
				progressWidth := progress * p.Width / p.Total
				percent := float64(progress) / float64(p.Total) * 100
				fmt.Printf("\r%s[%s%s]%d/%d %.2f%%",
					p.Title,
					strings.Repeat("=", progressWidth),
					strings.Repeat("_", p.Width-progressWidth),
					progress,
					p.Total,
					percent,
				)
				if progress == p.Total {
					fmt.Println()
					return
				}
			case <-ticker.C:
			}
		}
	}()
}

// Add 更新进度条
func (p *ProgressBar) Add(n int) {
	p.Lock()
	p.Progress = p.Progress + n
	p.ProgressChan <- p.Progress
	time.Sleep(time.Millisecond * 10)
	p.Unlock()
}
