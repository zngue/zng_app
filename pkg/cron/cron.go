package cron

type ICron interface {
	Run()
	Stop()
}
type ICronServer struct{}

func (ICronServer) Run() {
	panic("implement me ICronServer.Run")
}
func (ICronServer) Stop() {
	panic("implement me ICronServer.Stop")
}
