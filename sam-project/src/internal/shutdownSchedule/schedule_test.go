package schedule

import (
	"home_services_analyst/internal/config"
	"home_services_analyst/internal/repository"
	"home_services_analyst/internal/telegram"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/suite"
)

type ScheduleServiceTestSuite struct {
	suite.Suite
	schedule ShutdownSchedule
}

func TestScheduleServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduleServiceTestSuite))
}

func (s *ScheduleServiceTestSuite) SetupTest() {
	var telegramTestConfig config.TelegramConfig
	var repoConfig config.RepoConfig

	if err := cleanenv.ReadConfig("../../secrets/.dev.env", &telegramTestConfig); err != nil {
		panic(err)
	}
	if err := cleanenv.ReadConfig("../../secrets/.dev.env", &repoConfig); err != nil {
		panic(err)
	}

	telegram := telegram.NewTelegramService(&telegramTestConfig)
	repo := repository.NewRepository(&repoConfig)

	s.schedule = ShutdownSchedule{
		telegram: &telegram,
		repo:     &repo,
	}
}

func (s *ScheduleServiceTestSuite) TestLoadMetadata() {
	meta, err := s.schedule.repo.LoadMetadata()
	if err != nil {
		s.T().Error(err)
	}

	s.T().Log(meta)
}

func (s *ScheduleServiceTestSuite) TestHandleShutdownScheduleUpdate() {
	if err := s.schedule.HandleShutdownScheduleUpdate(); err != nil {
		s.Error(err)
	}
}

func (s *ScheduleServiceTestSuite) TestGetScheduleImage() {
	image, err := s.schedule.getScheduleImage()
	if err != nil {
		s.Error(err)
	}

	// expect image to be greater than 20k
	s.Greater(len(image), 20000)
}

func (s *ScheduleServiceTestSuite) TestSendNewImageInChannel() {
	image, err := s.schedule.getScheduleImage()
	if err != nil {
		s.Error(err)
	}

	if err := s.schedule.sendNewImageInChannel(&image); err != nil {
		s.Error(err)
	}
}

func (s *ScheduleServiceTestSuite) TestCalculateMd5() {
	image := make([]byte, 400)

	hash := s.schedule.calculateMd5(image)

	s.Equal("a75d7d422fd00bf31208b013e74d8394", hash)
}
