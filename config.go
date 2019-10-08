package main

import (
	"github.com/juju/loggo"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/env"
	"github.com/micro/go-micro/config/source/file"
	"slack-bot/workers"
	"strings"
)

type logger struct {
	Level string
}

type slackConfig struct {
	Token 			string `json:"token"`
	ExcludedUsers 	[]string `json:"excluded_user_ids"`
}

type Configuration struct{
	Loggers  map[string]logger  `json:"loggers"`
	Slack  	slackConfig      	`json:"slack"`
	Workers	[]workers.Config 	`json:"workers"`
}

func loggerConfigInitiator(conf *map[string]logger){
	var configLoggers []string
	for name, logger := range *conf {
		configLoggers = append(configLoggers, name+"="+logger.Level)
	}
	err := loggo.ConfigureLoggers(strings.Join(configLoggers,";"))
	if err != nil {}
}

func configConstructor() Configuration {
	var conf Configuration
	loadConfigErr := config.Load(
						env.NewSource(
							env.WithStrippedPrefix("SB"),
						),
						file.NewSource(
							file.WithPath("slack-bot.yaml"),
							file.WithPath("slack-bot-template.yaml"),
						),
					)
	if loadConfigErr != nil {	}
	syncConfigErr := config.Sync()
	if syncConfigErr != nil {	}
	parseConfigErr := config.Scan(&conf)
	if parseConfigErr != nil { }
	loggerConfigInitiator(&conf.Loggers)
	return conf
}