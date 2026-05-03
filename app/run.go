package app

// type App struct {
// 	httpSrv *http.Server
// 	cron    []cron.ICron
// 	Port    int32
// }

// func NewRouter(items []router.IApiService) (routes []router.IRouter) {
// 	for _, service := range items {
// 		runItems := service.Register()
// 		if len(runItems) > 0 {
// 			for _, runItem := range runItems {
// 				routes = append(routes, runItem)
// 			}
// 		}
// 	}
// 	return
// }

// func Abc()
// func NewApp(server *http.Server, cron []cron.ICron) *App {
// 	return &App{
// 		httpSrv: server,
// 		cron:    cron,
// 	}
// }
// func NewAppRunner(port int32, fn Fn) (err error) {
// 	var (
// 		cleanup func()
// 		run     *App
// 	)
// 	run, cleanup, err = fn()
// 	if err != nil {
// 		return
// 	}
// 	defer cleanup()
// 	fmt.Printf("http://127.0.0.1:%d\n", port)
// 	fmt.Printf("http://localhost:%d\n", port)
// 	err = run.Run()
// 	return
// }

// var ProviderSet = wire.NewSet(
// 	NewRouter,
// 	NewApp,
// )
