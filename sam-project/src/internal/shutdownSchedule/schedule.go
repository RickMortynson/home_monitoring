package schedule

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"home_services_analyst/internal/repository"
	"home_services_analyst/internal/telegram"
	"home_services_analyst/internal/utils"
	"io"
	"net/http"
)

const (
	SCHEDULE_IMAGE_URL    = "http://oblenergo.cv.ua/shutdowns/GPV.png"
	SCHEDULE_UPDATED_TEXT = "🔔 Розклад відключень оновлено"
)

type ShutdownSchedule struct {
	repo     *repository.Repository
	telegram *telegram.Telegram
}

func NewShutdownScheduleService(repo *repository.Repository, tg *telegram.Telegram) ShutdownSchedule {
	return ShutdownSchedule{
		repo:     repo,
		telegram: tg,
	}
}

func (s *ShutdownSchedule) HandleShutdownScheduleUpdate() error {
	metadata, err := s.repo.LoadMetadata()
	if err != nil {
		return err
	}
	utils.PrettyPrint(metadata)

	image, err := s.getScheduleImage()
	if err != nil {
		return err
	}

	hashedImage := s.calculateMd5(image)

	fmt.Printf("calculated new hashed Image against old one:\n %s\n%s \n",
		hashedImage,
		metadata.PowerScheduleMd5)

	if hashedImage != metadata.PowerScheduleMd5 {
		if err := s.sendNewImageInChannel(&image); err != nil {
			return err
		}

		s.repo.UpdateMetadata(repository.PostMetadata{
			Online:           metadata.Online,
			LastBadMessageId: metadata.LastBadMessageId,
			PowerScheduleMd5: hashedImage,
		})
	}

	return nil
}

func (s *ShutdownSchedule) sendNewImageInChannel(image *[]byte) error {
	if err := s.telegram.SendImageMessage(*image, SCHEDULE_UPDATED_TEXT); err != nil {
		return err
	}

	return nil
}

func (s *ShutdownSchedule) calculateMd5(image []byte) string {
	hash := md5.Sum(image)

	return hex.EncodeToString(hash[:])
}

func (s *ShutdownSchedule) getScheduleImage() ([]byte, error) {
	resp, err := http.Get(SCHEDULE_IMAGE_URL)
	if err != nil {
		return nil, err
	}

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return image, nil
}
