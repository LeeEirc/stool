package main

import (
	"flag"
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
		slog.Info("handle replay file: %s", f.AbsFilePath)

		session, err := apiClient.GetSessionById(f.ID)
		if err != nil {
			slog.Error("get session %s failed: %s", f.ID, err)
			continue
		}
		if session.ID == "" {
			slog.Warn("session %s not found", f.ID)
			continue
		}
		if !session.IsFinished {
			slog.Warn("session %s not finished", f.ID)
			continue
		}
		if !session.HasReplay {
			slog.Warn("session %s has replay", f.ID)
			continue
		}
		// 上传录像，并删除录像文件
		if err := apiClient.UploadReplay(f.ID, f.AbsFilePath); err != nil {
			slog.Error("upload replay %s failed: %s", f.ID, err)
			continue
		}
		if err := apiClient.FinishReply(f.ID); err != nil {
			slog.Error("finish replay %s failed: %s", f.ID, err)
			continue
		}
		slog.Info("handle replay file: %s success", f.AbsFilePath)
		if err1 := os.Remove(f.AbsFilePath); err1 != nil {
			slog.Error("delete replay file %s failed: %s", f.AbsFilePath, err1)
			continue
		}
		slog.Info("delete replay file %s success", f.AbsFilePath)
	}
}
