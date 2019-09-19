package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"TlsEasyDarwin/models"
	"TlsEasyDarwin/penggy/EasyGoLib/utils"
	"TlsEasyDarwin/penggy/service"
	"TlsEasyDarwin/routers"
	"TlsEasyDarwin/rtsp"
	"github.com/common-nighthawk/go-figure"
)

type program struct {
	httpPort   int
	httpServer *http.Server
	rtspPort   int
	rtspServer *rtsp.Server
	// 使用tls时才有rtspsServer
	rtspsPort   int
	rtspsServer *rtsp.Server
}

func (p *program) StopHTTP() (err error) {
	if p.httpServer == nil {
		err = fmt.Errorf("HTTP Server Not Found")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = p.httpServer.Shutdown(ctx); err != nil {
		return
	}
	return
}

func (p *program) StartHTTP() (err error) {
	p.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", p.httpPort),
		Handler:           routers.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	link := fmt.Sprintf("http://%s:%d", utils.LocalIP(), p.httpPort)
	log.Println("http server start -->", link)
	go func() {
		if err := p.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("start http server error", err)
		}
		log.Println("http server end")
	}()
	return
}

func (p *program) StartRTSP() (err error) {
	if p.rtspServer == nil {
		err = fmt.Errorf("RTSP Server Not Found")
		return
	}
	sport := ""
	if p.rtspPort != 554 {
		sport = fmt.Sprintf(":%d", p.rtspPort)
	}
	link := fmt.Sprintf("rtsp://%s%s", utils.LocalIP(), sport)
	log.Println("rtsp server start -->", link)
	go func() {
		if err := p.rtspServer.Start(); err != nil {
			log.Println("start rtsp server error", err)
		}
		log.Println("rtsp server end")
	}()
	return
}

func (p *program) StopRTSP() (err error) {
	if p.rtspServer == nil {
		err = fmt.Errorf("RTSP Server Not Found")
		return
	}
	p.rtspServer.Stop()
	return
}

func (p *program) StartRTSPS() (err error) {
	err = nil
	if p.rtspsServer == nil {
		// err = fmt.Errorf("RTSPS Server Not Found")
		return nil
	}
	sport := ""
	if p.rtspsPort != 554 {
		sport = fmt.Sprintf(":%d", p.rtspsPort)
	}
	link := fmt.Sprintf("rtsps://%s%s", utils.LocalIP(), sport)
	log.Println("rtsps server start -->", link)
	go func() {
		if err := p.rtspsServer.Start(); err != nil {
			log.Println("start rtsps server error", err)
		}
		log.Println("rtsps server end")
	}()
	return
}

func (p *program) StopRTSPS() (err error) {
	err = nil
	if p.rtspsServer == nil {
		// err = fmt.Errorf("RTSPS Server Not Found")
		return nil
	}
	p.rtspsServer.Stop()
	return
}

func (p *program) Start(s service.Service) (err error) {
	log.Println("********** START **********")
	if utils.IsPortInUse(p.httpPort) {
		err = fmt.Errorf("HTTP port[%d] In Use", p.httpPort)
		return
	}
	if utils.IsPortInUse(p.rtspPort) {
		err = fmt.Errorf("RTSP port[%d] In Use", p.rtspPort)
		return
	}
	// 存在rtspsServer时
	if p.rtspsServer != nil {
		if utils.IsPortInUse(p.rtspsPort) {
			err = fmt.Errorf("RTSPS port[%d] In Use", p.rtspsPort)
			return
		}
	}

	err = models.Init()
	if err != nil {
		return
	}
	err = routers.Init()
	if err != nil {
		return
	}
	// tls
	p.StartRTSPS()
	p.StartRTSP()
	p.StartHTTP()
	if !utils.Debug {
		log.Println("log files -->", utils.LogDir())
		log.SetOutput(utils.GetLogWriter())
	}
	go func() {
		for range routers.API.RestartChan {
			p.StopHTTP()
			p.StopRTSP()
			p.StopRTSPS()
			utils.ReloadConf()
			p.StopRTSPS()
			p.StartRTSP()
			p.StartHTTP()
		}
	}()
	return
}

func (p *program) Stop(s service.Service) (err error) {
	defer log.Println("********** STOP **********")
	defer utils.CloseLogWriter()
	p.StopHTTP()
	p.StopRTSP()
	p.StopRTSPS()
	models.Close()
	return
}


/*
使用tls，要推流至rtsps
ffmpeg -re -i test.mp4 -rtsp_transport tcp -c copy -f rtsp rtsps://127.0.0.1:8443/test.mp4
播放
ffplay rtsps://10.5.15.57:8443/test.mp4
不使用tls，要推流至rtsp
ffmpeg -re -i test.mp4 -rtsp_transport tcp -c copy -f rtsp rtsp://127.0.0.1:8554/test.mp4
播放
ffplay rtsp://10.5.15.57:8554/test.mp4

可执行文件名必须为EasyDarwin(不区分大小写),因为根据此读取配置文件easydarwin.ini
// TODO
// 1. 配置文件更改，证书、测试相关
// 2. 日志更改
// 3. README更改
*/
func main() {

	log.SetPrefix("[EasyDarwin] ")
	log.SetFlags(log.LstdFlags)
	if utils.Debug {
		log.SetFlags(log.Lshortfile | log.LstdFlags)
	}
	// 编译release版本，日志写入文件
	// go build -tags release -o EasyDarwin main.go
	// 日志写入logs目录内文件
	log.SetOutput(utils.GetLogWriter())

	sec := utils.Conf().Section("service")
	svcConfig := &service.Config{
		Name:        sec.Key("name").MustString("EasyDarwin_Service"),
		DisplayName: sec.Key("display_name").MustString("EasyDarwin_Service"),
		Description: sec.Key("description").MustString("EasyDarwin_Service"),
	}

	httpPort := utils.Conf().Section("http").Key("port").MustInt(10008)

	var rtspsSrv *rtsp.Server
	var rtspsPort int
	// 配置文件内有tls时,启用tls; 没有时tls=false
	_, e := utils.Conf().GetSection("tls")
	if e != nil {
		rtspsSrv = nil
		rtspsPort = 443
	} else {
		rtspsPort = utils.Conf().Section("tls").Key("port").MustInt(443)
		rtspsSrv = rtsp.NewServer(rtspsPort, true)
	}

	rtspServer := rtsp.GetServer()
	p := &program{
		httpPort:   httpPort,
		rtspPort:   rtspServer.TCPPort,
		rtspServer: rtspServer,
		// 使用tls,rtsps server;
		rtspsPort: rtspsPort,
		rtspsServer: rtspsSrv,
	}
	var s, err = service.New(p, svcConfig)
	if err != nil {
		log.Println(err)
		utils.PauseExit()
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" || os.Args[1] == "stop" {
			figure.NewFigure("EasyDarwin", "", false).Print()
		}
		log.Println(svcConfig.Name, os.Args[1], "...")
		if err = service.Control(s, os.Args[1]); err != nil {
			log.Println(err)
			utils.PauseExit()
		}
		log.Println(svcConfig.Name, os.Args[1], "ok")
		return
	}
	figure.NewFigure("EasyDarwin", "", false).Print()
	if err = s.Run(); err != nil {
		log.Println(err)
		utils.PauseExit()
	}
}
