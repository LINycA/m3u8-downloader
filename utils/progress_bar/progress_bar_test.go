package progressbar_test

import (
	progressbar "m3u8-downloader/utils/progress_bar"
	"testing"
)

func TestProgressBar(t *testing.T) {
	total := 150
	progressBar := progressbar.NewProgressBar("test", 150)
	for range total {
		progressBar.Add(1)
	}
}
