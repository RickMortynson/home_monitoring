package main

import (
	"fmt"
	"home_services_analyst/internal/config"
	"home_services_analyst/internal/heartbeat"
	"home_services_analyst/internal/repository"
	"home_services_analyst/internal/telegram"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler() {
	var tgConfig config.TelegramConfig
	var heartbeatConfig config.HeartbeatConfig
	var repoConfig config.RepoConfig

	if err := config.LoadConfig(&tgConfig, &heartbeatConfig, &repoConfig); err != nil {
		fmt.Fprintf(os.Stderr, "❌❌❌%v❌❌❌", err)
		panic(err)
	}

	fmt.Println("configs:", tgConfig, heartbeatConfig, repoConfig)

	repo := repository.NewRepository(&repoConfig)
	telegram := telegram.NewTelegramService(&tgConfig)

	heartbeatService := heartbeat.NewHeartbeatService(&heartbeatConfig, &repo, &telegram)

	if err := heartbeatService.HandleHeartbeats(); err != nil {
		fmt.Fprintf(os.Stderr, "❌❌❌%v❌❌❌", err)
		panic(err)
	}
}
