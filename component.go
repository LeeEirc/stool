package main

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

type Component struct {
	Name       string
	ReplayPath string
	cfg        *Config
}

func (c *Component) ScanSessionReplays() []ReplayFile {
	files := make([]ReplayFile, 0, 100)
	nowDate := time.Now().Format("2006-01-02")
	_ = filepath.Walk(c.ReplayPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// 必须是 .gz 文件
		// 必须是 sid 文件
		// 必须目录日期小于当前日期
		// 必须是 session id
		if !strings.HasSuffix(info.Name(), ".gz") {
			return nil
		}
		dateStr, ok := ParseDateFromPath(path)
		if !ok {
			return nil
		}
		if dateStr == nowDate {
			return nil
		}
		sid, ok := ParseSessionID(path)
		if !ok {
			return nil
		}
		version, ok := ParseReplayVersion(info.Name())
		if !ok {
			return nil
		}
		rf := ReplayFile{
			ID:          sid,
			TargetDate:  dateStr,
			AbsFilePath: path,
			Version:     version,
		}
		files = append(files, rf)
		return nil
	})
	return files
}
