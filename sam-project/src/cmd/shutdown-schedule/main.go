package main

import (
	"fmt"
	"home_services_analyst/internal/config"
	"home_services_analyst/internal/repository"
	schedule "home_services_analyst/internal/shutdownSchedule"
	"home_services_analyst/internal/telegram"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler() {
	fmt.Println("USER", os.Getenv("USER"))

	var tgConfig config.TelegramConfig
	var repoConfig config.RepoConfig

	if err := config.LoadConfig(&tgConfig, &repoConfig); err != nil {
		fmt.Fprintf(os.Stderr, "❌❌❌%v❌❌❌", err)
		panic(err)
	}

	repo := repository.NewRepository(&repoConfig)
	telegram := telegram.NewTelegramService(&tgConfig)

	scheduleService := schedule.NewShutdownScheduleService(&repo, &telegram)

	if err := scheduleService.HandleShutdownScheduleUpdate(); err != nil {
		fmt.Fprintf(os.Stderr, "❌❌❌%v❌❌❌", err)
		panic(err)
	}
}
