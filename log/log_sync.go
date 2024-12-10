package log

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type SaveLog struct {
	Url string
}
type Data struct {
	Level string `json:"level"`
}

func (s *SaveLog) Write(p []byte) (n int, err error) {
	var rs Data
	err = json.Unmarshal(p, &rs)
	if err != nil {
		return
	}
	if rs.Level != "error" {
		return
	}
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
