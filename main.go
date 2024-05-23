package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"stool/jms-sdk-go/httplib"
	"stool/jms-sdk-go/service"
)

var cfgPath = "config.yml"

func init() {
	flag.StringVar(&cfgPath, "c", "config.yml", "config file path")
}

func main() {
	flag.Parse()
	cfg := LoadConfig(cfgPath)
	components := make([]Component, 0, 10)
	for _, path := range cfg.ReplayPaths {
		c := Component{
			Name:       path,
			ReplayPath: path,
			cfg:        cfg,
		}
		components = append(components, c)
	}
	sign := httplib.PrivateTokenAuth{Token: cfg.PrivateToken}
	apiClient, err := service.NewAuthJMService(service.JMSCoreHost(
		cfg.CoreHost), service.JMSTimeOut(30*time.Second),
		service.JMSAuthSign(&sign),
	)
	if err != nil {
		slog.Error("create jms service failed: " + err.Error())
		os.Exit(1)
	}
	handleComponents(components, apiClient)
	slog.Info("all done")
}

func handleComponents(components []Component, apiClient *service.JMService) {
	var wg sync.WaitGroup

	for i := range components {
		wg.Add(1)
		go handleComponent(&wg, &components[i], apiClient)
	}
	wg.Wait()

}

func handleComponent(wg *sync.WaitGroup, c *Component, apiClient *service.JMService) {
	defer wg.Done()
	files := c.ScanSessionReplays()
	for _, f := range files {
		msg := fmt.Sprintf("Start handle replay file: %s", f.AbsFilePath)
		slog.Info(msg)

		session, err := apiClient.GetSessionById(f.ID)
		if err != nil {
			msg = fmt.Sprintf("Get session %s failed: %s", f.ID, err)
			slog.Error(msg)
			continue
		}
		if session.ID == "" {
			slog.Error("Not found session " + f.ID)
			continue
		}
		if !session.IsFinished {
			slog.Error("Not finished session " + f.ID)
			continue
		}
		if session.HasReplay {
			if !c.cfg.OverWriteReplay {
				msg = fmt.Sprintf("Session %s alreay have replay", f.ID)
				slog.Error(msg)
				continue
			}
			msg = fmt.Sprintf("Session %s alreay have replay, try to overwrite it", f.ID)
			slog.Info(msg)
		}
		msg = fmt.Sprintf("Uploading session replay file %s", f.AbsFilePath)
		// 上传录像，并删除录像文件
		if err1 := apiClient.UploadReplay(f.ID, f.AbsFilePath); err1 != nil {
			msg = fmt.Sprintf("Uploading replay %s failed: %s", f.ID, err1)
			slog.Error(msg)
			continue
		}
		msg = fmt.Sprintf("Upload replay file %s success", f.AbsFilePath)
		slog.Info(msg)
		if err2 := apiClient.FinishReply(f.ID); err2 != nil {
			errMsg := fmt.Sprintf("Finish replay %s failed: %s", f.ID, err2)
			slog.Error(errMsg)
			continue
		}
		msg = fmt.Sprintf("Finish session %s success", f.ID)
		slog.Info(msg)
		if err1 := os.Remove(f.AbsFilePath); err1 != nil {
			errMsg := fmt.Sprintf("Delete replay file %s failed: %s", f.AbsFilePath, err1)
			slog.Error(errMsg)
			continue
		}
		msg = fmt.Sprintf("Delete replay file %s success", f.AbsFilePath)
		slog.Info(msg)
		msg = fmt.Sprintf("Finish handle replay file %s", f.AbsFilePath)
		slog.Info(msg)
	}
}
