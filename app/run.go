package app

import (
	"fmt"
	"github.com/google/wire"
	"github.com/zngue/zng_app/pkg/cron"
	"github.com/zngue/zng_app/pkg/router"
	"net/http"
)

type App struct {
	httpSrv *http.Server
	cron    []cron.ICron
	routers []router.IRouter
	Port    int32
}

func NewRouter(items []router.IApiService) (routes []router.IRouter) {
	for _, service := range items {
		runItems := service.Register()
		if len(runItems) > 0 {
			for _, runItem := range runItems {
				routes = append(routes, runItem)
			}
		}
	}
	return
}
func NewApp(server *http.Server, routers []router.IRouter, cron []cron.ICron) *App {
	return &App{
		httpSrv: server,
		routers: routers,
		cron:    cron,
	}
}
func NewAppRunner(port int32, fn Fn) (err error) {
	var (
		cleanup func()
		run     *App
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

var ProviderSet = wire.NewSet(
	NewRouter,
	NewApp,
)
