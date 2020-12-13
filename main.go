/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-08 14:52
* Description:
*****************************************************************/

package main

import (
	"context"
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/x/utils/xconfig"
	"github.com/go-xe2/xthrift/lib/go/xthrift"
	"github.com/go-xe2/xthrift/rpcRouter"
	"log"
	"net"
	"os"
)

var server *rpcRouter.TRouterServer
var enableRpcServer = false
var rpcServerAddr string
var rpcSvr *xthrift.TXServer
//
//func appInit(ctx context.Context) ([]net.Listener, error)  {
//	options, ok := ctx.Value("options").(*rpcRouter.TSvrOptions)
//	if !ok {
//		return nil, errors.New("参数options为nil")
//	}
//	httpLst, err := net.Listen("tcp", options.HttpAddr)
//	if err != nil {
//		return nil, err
//	}
//	rpcLst, err := net.Listen("tcp", options.RouterAddr)
//	if err != nil {
//		return nil, err
//	}
//	return []net.Listener{httpLst, rpcLst}, nil
//}
//
//func appRun(ctx context.Context, listeners []net.Listener) error  {
//	options := ctx.Value("options").(*rpcRouter.TSvrOptions)
//	server = rpcRouter.NewServer(options, listeners[0], listeners[1])
//
//	if enableRpcServer {
//		processorFac := rpcRouter.NewRouterProcessorFactory(server)
//
//		trans, err := xthrift.NewServer(rpcServerAddr)
//		if err != nil {
//			panic(err)
//		}
//		rpcSvr = trans
//		rpcSvr.SetProcessorFac(processorFac)
//		go func() {
//			xlog.Info("rpc服务端口:", rpcServerAddr)
//			if err := rpcSvr.Serve(); err != nil {
//				xlog.Error(err)
//			}
//		}()
//	}
//	xlog.Info("协议存放目录:", options.PDLPath)
//	xlog.Info("服务器地址存放目录:", options.HostPath)
//	xlog.Info("路由服务监听地址:", options.RouterAddr)
//	xlog.Info("http服务监听地址:", options.HttpAddr)
//	return server.Serve()
//}
//
//func appStop(ctx context.Context) {
//	if rpcSvr != nil {
//		if err := rpcSvr.Stop(); err != nil {
//			xlog.Error(err)
//		}
//	}
//	if server != nil {
//		if err := server.Stop(); err != nil {
//			xlog.Error(err)
//		}
//	}
//}


func main() {
	//loader := hotLoad.NewHotLoadService("rpc-router.pid")
	c := xconfig.Config()
	c.SetFileName("config.yaml")

	mp := c.GetMap("options")

	if v, ok := mp["enableRpcServer"]; ok {
		enableRpcServer = t.Bool(v)
	}
	if v, ok := mp["rpcAddr"]; ok {
		rpcServerAddr = t.String(v)
	}

	if !xfile.Exists("./logs") {
		_ = xfile.Mkdir("./logs")
	}
	logFile, logErr := os.OpenFile(xfile.Join("./logs", "rpc-router.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("测试日志输出文件出错:", logErr)
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
	}
	if err := xlog.SetPath("./logs"); err != nil {
		fmt.Println("设置日志输出目录出错:", err)
	}

	options, err := rpcRouter.NewSvrOptionsFromMap(mp)
	if err != nil {
		xlog.Info(err)
		return
	}
	cxt := context.Background()
	cxt = context.WithValue(cxt, "options", options)

	httpLst, err := net.Listen("tcp", options.HttpAddr)
	if err != nil {
		panic(err)
	}
	rpcLst, err := net.Listen("tcp", options.RouterAddr)
	if err != nil {
		panic(err)
	}

	server = rpcRouter.NewServer(options, httpLst, rpcLst)

	if enableRpcServer {
		processorFac := rpcRouter.NewRouterProcessorFactory(server)

		trans, err := xthrift.NewServer(rpcServerAddr)
		if err != nil {
			panic(err)
		}
		rpcSvr = trans
		rpcSvr.SetProcessorFac(processorFac)
		go func() {
			xlog.Info("rpc服务端口:", rpcServerAddr)
			if err := rpcSvr.Serve(); err != nil {
				xlog.Error(err)
			}
		}()
	}
	xlog.Info("协议存放目录:", options.PDLPath)
	xlog.Info("服务器地址存放目录:", options.HostPath)
	xlog.Info("路由服务监听地址:", options.RouterAddr)
	xlog.Info("http服务监听地址:", options.HttpAddr)
	if err := server.Serve(); err != nil {
		xlog.Debug(err)
	}
	//err = loader.Load(cxt, appInit, appRun, appStop)
	//if err != nil {
	//	xlog.Error(err)
	//}
}