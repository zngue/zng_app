package main

import (
	"fmt"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/workflow"
	"github.com/lithammer/shortuuid/v3"
)

func main() {
	const clientUrl = "http://172.16.26.213:16666/v1/dtm"
	dtmServer := "http://172.29.123.236:36789/api/dtmsvr"
	var req = make(map[string]any)
	req["name"] = "acb"
	saga := dtmcli.NewSaga(dtmServer, shortuuid.New())
	saga.Add(clientUrl+"/transOut", clientUrl+"/transOutCompensate", req)
	saga.Add(clientUrl+"/transIn", clientUrl+"/transInCompensate", req)
	err := saga.Submit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("over")
}

func Grpc() {
	dtmServer := "http://172.29.123.236:36789/api/dtmsvr"
	workflow.InitGrpc(dtmServer, busi.BusiGrpc, s)
	wfName := "workflow-grpc"
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var req busi.BusiReq
		err := proto.Unmarshal(data, &req)
		logger.FatalIfError(err)
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := busiCli.TransOutRevert(wf.Context, &req)
			return err
		})
		_, err = busiCli.TransOut(wf.Context, &req)
		if err != nil {
			return err
		}

		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			_, err := busiCli.TransInRevert(wf.Context, &req)
			return err
		})
		_, err = busiCli.TransIn(wf.Context, &req)
		return err
		return nil
	})

}
