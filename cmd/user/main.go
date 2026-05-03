package main

import (
	"github.com/zngue/zng_app/app"
	"github.com/zngue/zng_app/app/config/nacos"
)

func main() {

	var err = nacos.NewCenterOptions()
	app.New()
	if err != nil {
		panic(err)
	}
}
