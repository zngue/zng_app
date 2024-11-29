package app

type IRun interface {
	Run()
}
type IRunServer struct {
}

func (IRunServer) Run() {
	panic("Unimplemented IRunServer.Run")
}
