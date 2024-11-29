package app

import "context"

type IApp interface {
	Run() (err error)
	Stop(ctx context.Context) error
}
type IAppServer struct{}

func (IAppServer) Run() (err error) {
	panic("implement me IAppServer.Run")
}
func (IAppServer) Stop(ctx context.Context) error {
	panic("implement me IAppServer.Stop")
}
