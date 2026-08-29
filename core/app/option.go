package app

import "time"

type Option func(*App)

func WithAppName(name string) Option {
	return func(a *App) {
		a.appName = name
	}
}

func WithServer(s Server) Option {
	return func(a *App) {
		a.servers = append(a.servers, s)
	}
}

func WithCloser(c Closer) Option {
	return func(a *App) {
		a.closers = append(a.closers, c)
	}
}

func WithShutdownTimeout(d time.Duration) Option {
	return func(a *App) {
		a.timeout = d
	}
}
