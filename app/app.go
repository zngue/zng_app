package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	httpSrv *http.Server
	apps    []IApp
	routers []IRouter
}

func NewHttpServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}
func NewApp(server *http.Server, routers []IRouter, apps []IApp) *App {
	return &App{
		httpSrv: server,
		routers: routers,
		apps:    apps,
	}
}
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
		for _, app := range a.apps {
			if appErr := app.Run(); appErr != nil {
				panic(appErr)
			}
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
		if len(a.apps) > 0 {
			// 关闭其他服务
			for _, app := range a.apps {
				if err := app.Stop(ctx); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()
	return nil
}

func NewAppRunner(port int32, fn Fn) (err error) {
	var (
		run     *App
		cleanup func()
	)
	run, cleanup, err = fn()
	if err != nil {
		return
	}
	defer cleanup()
	fmt.Println(fmt.Sprintf("http://127.0.0.1:%d", port))
	fmt.Println(fmt.Sprintf("http://localhost:%d", port))
	err = run.Run()
	return
}
