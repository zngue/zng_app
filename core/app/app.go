package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	zng_app "github.com/zngue/zng_app"
	"golang.org/x/sync/errgroup"
)

type Server interface {
	Start() error
	Stop() error
}

type Closer func() error

type App struct {
	appName string
	servers []Server
	closers []Closer
	timeout time.Duration

	mu      sync.Mutex
	started bool
}

func New(opts ...Option) *App {
	a := &App{
		appName: zng_app.AppName,
		timeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(a)
	}
	if a.appName != "" {
		zng_app.SetAppName(a.appName)
	}
	return a
}

func (a *App) AppName() string {
	return a.appName
}

func (a *App) AddCloser(c Closer) {
	a.closers = append(a.closers, c)
}

func (a *App) Run() (err error) {
	a.mu.Lock()
	if a.started {
		a.mu.Unlock()
		return errors.New("app: already started")
	}
	a.started = true
	a.mu.Unlock()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg errgroup.Group
	for _, s := range a.servers {
		wg.Go(s.Start)
	}

	select {
	case <-ctx.Done():
		fmt.Printf("\n[%s] app: shutting down...\n", a.appName)
	case err = <-gWait(&wg):
		return fmt.Errorf("app: server error: %w", err)
	}

	if shutdownErr := a.Stop(); shutdownErr != nil {
		return fmt.Errorf("app: shutdown error: %w", shutdownErr)
	}
	return
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	var errs []error
	for i := len(a.servers) - 1; i >= 0; i-- {
		if err := a.servers[i].Stop(); err != nil {
			errs = append(errs, err)
		}
	}
	for i := len(a.closers) - 1; i >= 0; i-- {
		if err := a.closers[i](); err != nil {
			errs = append(errs, err)
		}
	}
	<-ctx.Done()
	return errors.Join(errs...)
}

func MustRun(a *App) {
	if err := a.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] app: %v\n", a.appName, err)
		os.Exit(1)
	}
}

func gWait(g *errgroup.Group) <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- g.Wait()
	}()
	return ch
}
