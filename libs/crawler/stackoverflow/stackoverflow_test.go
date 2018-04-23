package stackoverflow

import (
	"fmt"
	"testing"

	"github.com/fengyfei/gu/libs/crawler"
)

func TestNewXteam(t *testing.T) {
	var (
		dataCh   = make(chan crawler.Data)
		finishCh = make(chan struct{})
	)
	c := NewStackOverFlow(dataCh, finishCh)
	go func() {
		err := crawler.StartCrawler(c)
		if err != nil {
			panic(err)
		}
	}()
	for {
		select {
		case data := <-dataCh:
			if data != nil {
				fmt.Println(data.(*Blog).Date)
			}
		case <-finishCh:
			return
		}
	}
}
