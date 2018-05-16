package main

import (
	"net"
	"net/http"
	"strconv"

	"github.com/clzhan/srlive3/conf"
	"github.com/clzhan/srlive3/httpserver"
	"github.com/clzhan/srlive3/log"
	"github.com/clzhan/srlive3/protocol/rtmp"
	"github.com/clzhan/srlive3/utils"
)

//远程获取pprof数据
func InitPprof() {
	//获取本机ip
	//rtmpAddr := fmt.Sprintf("%s:%d", network.GetLocalIpAddress(),6399)
	//
	//str ,_ := network.IntranetIP()
	//log.Info("local ip: ",str)

	go func() {
		//http://10.10.6.162:6399/debug/pprof
		pprofAddress := util.GetLocalIp()
		pprofAddress += ":"
		pprofAddress += strconv.Itoa(6399)

		log.Info(http.ListenAndServe(pprofAddress, nil))
	}()

}

func startHttpServer() error {
	var httpServerListen net.Listener
	var err error

	HttpFlvAddress := util.GetLocalIp()
	HttpFlvAddress += ":"
	HttpFlvAddress += strconv.Itoa(conf.AppConf.HttpPort)

	log.Info("HTTP-Server listen On", HttpFlvAddress)
	httpServerListen, err = net.Listen("tcp", HttpFlvAddress)

	if err != nil {
		log.Error(err)
		return err
	}

	httpServer := httpserver.NewHttpServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HTTP server panic: ", r)
			}
		}()
		log.Info("HttpServer listen On", HttpFlvAddress)
		httpServer.Serve(httpServerListen)
	}()
	return err
}

func main() {

	conf.Init()
	log.Init()

	err := startHttpServer()
	log.Info("ListenAndServerHttpServer error :", err)

	RtmpAddress := util.GetLocalIp()
	RtmpAddress += ":"
	RtmpAddress += strconv.Itoa(conf.AppConf.RtmpPort)

	err = rtmp.ListenAndServer(RtmpAddress)
	if err != nil {
		panic(err)
	}

	log.Debug("rtmp ListenAndServer :", RtmpAddress)


	//rtmp.ConnectPull("rtmp://10.10.6.39:1935/live/movie")

	// do event loop
	select {}

}
