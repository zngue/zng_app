package app

import (
	"log"
	"time"

	"github.com/zngue/zng_app/app/server"
	"golang.org/x/sync/errgroup"
)

type App struct {
	servers []server.Server
	timeout time.Duration
}

type Option func(*App)

func WithServer(s server.Server) Option {
	return func(a *App) {
		a.servers = append(a.servers, s)
	}
}

func WithShutdownTimeout(d time.Duration) Option {
	return func(a *App) {

		a.timeout = d
	}
}

func New(opts ...Option) *App {
	a := &App{
		timeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *App) Run() error {
	var wg errgroup.Group
	for _, srv := range a.servers {
		wg.Go(func() error {
			return srv.Start()
		})
	}
	if err := wg.Wait(); err != nil {
		log.Printf("server error: %v", err)
	}
	log.Println("app exited")
	return nil
}
func NewAppRunner(opts ...Option) (err error) {
	app := New(opts...)
	err = app.Run()
	return
}
