package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Fn func() (*App, func(), error)

func (a *App) Run() (err error) {
	go func() {
		httpErr := a.httpSrv.ListenAndServe()
		if httpErr != nil && !errors.Is(httpErr, http.ErrServerClosed) {
			panic(httpErr) // 修复：这里应该是httpErr而不是err
		}
	}()
	go func() {
		for _, app := range a.cron {
			app.Run()
		}
	}()
	log.Printf("start app running")
	quit := make(chan os.Signal, 1) // 修复：添加缓冲区大小为1
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("shutdown app")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// 关闭应用
	if err = a.Stop(ctx); err != nil {
		panic(err)
	}
	return
}
func (a *App) Stop(ctx context.Context) error {
	// 修复：移除goroutine，直接执行关闭操作
	if err := a.httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	if len(a.cron) > 0 {
		// 关闭其他服务
		for _, app := range a.cron {
			app.Stop()
		}
	}
	return nil
}
