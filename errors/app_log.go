package errors

import (
	"github.com/zngue/zng_app"
	"github.com/zngue/zng_app/log"
)

func LogS(err error) {
	if zng_app.SyncLogger { // api接口错误自动记录日志
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("panic recover err:%v", r)
				}
			}()
			log.Error(GetStackTrace(err))
		}()
	}

}
