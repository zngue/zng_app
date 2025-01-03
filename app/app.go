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
	for _, r := range a.routers {
		r.Router()
	}
	go func() {
		httpErr := a.httpSrv.ListenAndServe()
		if httpErr != nil && !errors.Is(httpErr, http.ErrServerClosed) {
			panic(err)
		}
	}()
	go func() {
		for _, app := range a.cron {
			app.Run()
		}
	}()
	log.Printf("start app running")
	quit := make(chan os.Signal)
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
	go func() {
		if err := a.httpSrv.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	go func() {
		if len(a.cron) > 0 {
			// 关闭其他服务
			for _, app := range a.cron {
				app.Stop()
			}
		}
	}()
	return nil
}
