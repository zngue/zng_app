package nacos

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func Test_Nacos(t *testing.T) {

	err := NewCenterOptions()
	if err != nil {
		t.Error(err)
	}
	for {
		fmt.Println("111")
		time.Sleep(time.Second * 2)
	}

}
func Test_Bc(t *testing.T) {
	cli, err := NewNamingClient("42.193.55.92", "develop")
	if err != nil {
		t.Error(err)
	}
	if err := waitReady(cli, 10*time.Second); err != nil {
		log.Fatal("nacos client init failed:", err)
	}

	_, err = cli.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: "user-service",
		Ip:          "127.0.0.1",
		Port:        9001,
		Ephemeral:   true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Registered successfully!")

}
