package log

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type SaveLog struct {
	Url string
}

func (s *SaveLog) Write(p []byte) (n int, err error) {
	client := resty.New()
	req := client.R()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		//http://localhost:16666/v1/log/create
		_, err = req.SetBody(string(p)).Post(s.Url)
		if err != nil {
			return
		}
	}()
	return
}
func NewLogSave(url string) *SaveLog {
	l := new(SaveLog)
	l.Url = url
	return l
}
