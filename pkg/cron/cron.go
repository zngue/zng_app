package cron

type ICron interface {
	Run()
	Stop()
}
type IAppServer struct{}

func (IAppServer) Run() {
	panic("implement me IAppServer.Run")
}
func (IAppServer) Stop() {
	panic("implement me IAppServer.Stop")
}
