package app

type IRouter interface {
	Router()
}

type IRouterServer struct {
}

func (IRouterServer) Router() {
	panic("Unimplemented IRouterServer.Router")
}
