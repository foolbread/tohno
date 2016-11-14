/*
@author: foolbread
@time: 2016-11-14
@file:tohno/config/config.go
*/
package config

import (
	"github.com/foolbread/fbcommon/config"
	"github.com/foolbread/fbcommon/golog"
)

func InitConfig(file string) {
	golog.Info("tohno config initing...")
	conf, err := config.LoadConfigByFile(file)
	if err != nil {
		golog.Critical(err)
	}

	var section string
	var str string

	section = "file"
	str = conf.MustString(section, "filepath", "")
	golog.Info(section, "filepath:", str)
	g_conf.setFilePath(str)

	str = conf.MustString(section, "filebackpath", "")
	golog.Info(section, "filebackpath:", str)
	g_conf.setFileBackupPath(str)
}

func GetConfig() *tohnoConfig {
	return g_conf
}

var g_conf *tohnoConfig = new(tohnoConfig)

type tohnoConfig struct {
	filePath       string
	fileBackupPath string
}

func (c *tohnoConfig) setFilePath(s string) {
	c.filePath = s
}

func (c *tohnoConfig) setFileBackupPath(s string) {
	c.fileBackupPath = s
}

func (c *tohnoConfig) GetFilePath() string {
	return c.filePath
}

func (c *tohnoConfig) GetFileBackupPath() string {
	return c.fileBackupPath
}
