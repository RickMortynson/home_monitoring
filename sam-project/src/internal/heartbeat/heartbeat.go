package heartbeat

import (
	"encoding/json"
	"fmt"
	"home_services_analyst/internal/config"
	"home_services_analyst/internal/repository"
	"home_services_analyst/internal/telegram"

	"io"
	"net/http"
	"time"
)

const (
	GOOD_MESSAGE = "З'явилось світло 🎉"
	BAD_MESSAGE  = "Нема світла 🫠"
)

type HeartBeatResponse struct {
	Heartbeats []ReceivedHeartbeat
}

type ReceivedHeartbeat struct {
	Name     string
	Warning  int
	Error    int
	Age      int
	Status   string
	LastBeat time.Time
}

type Heartbeat struct {
	config   *config.HeartbeatConfig
	repo     *repository.Repository
	telegram *telegram.Telegram

	metadata         repository.Metadata
	shouldUpdateRepo bool
}

func NewHeartbeatService(cfg *config.HeartbeatConfig, repo *repository.Repository, tg *telegram.Telegram) Heartbeat {
	return Heartbeat{
		config:   cfg,
		repo:     repo,
		telegram: tg,
	}
}

func (h *Heartbeat) HandleHeartbeats() error {
	if err := h.repo.LoadMetadata(&h.metadata); err != nil {
		return err
	}

	prettyPrint(h.metadata)

	lastGoodHeartbeat, err := h.getLastGoodHeartbeatTimestamp()
	if err != nil {
		return err
	}

	fmt.Println("Time since last heartbeat:", time.Since(lastGoodHeartbeat))

	lastHeartbeatInOkTime := time.Since(lastGoodHeartbeat) < time.Duration(h.config.CheckTimeoutSec*int(time.Second))

	fmt.Println("h.config.CheckTimeoutSec", time.Duration(h.config.CheckTimeoutSec))
	fmt.Println("lastHeartbeatInOkTime", lastHeartbeatInOkTime)
	fmt.Println("h.metadata.Online", h.metadata.Online)

	if lastHeartbeatInOkTime {
		if !h.metadata.Online {
			h.updateMetaOnlineStatus(true)

			fmt.Println("🚨🚨🚨 SEND GOOD MESSAGE 🚨🚨🚨")

			_, err = h.telegram.SendMessage(GOOD_MESSAGE)
			if err != nil {
				return err
			}
		}
	} else {
		if h.metadata.Online {
			h.updateMetaOnlineStatus(false)

			fmt.Println("🚨🚨🚨 SEND BAD MESSAGE 🚨🚨🚨")

			resp, err := h.telegram.SendMessage(BAD_MESSAGE)
			if err != nil {
				return err
			}

			h.updateLastMessageId(resp.Result.MessageId)
		} else {
			newMessage := fmt.Sprintf("%s\n\n%s", BAD_MESSAGE, fmt.Sprintf("%s%s", "Світла не було ", secondsToReadableUkrainian(int(time.Since(lastGoodHeartbeat).Seconds()))))
			h.telegram.UpdateTextMessage(h.metadata.LastMessageId, newMessage)
		}
	}

	if h.shouldUpdateRepo {
		if err := h.repo.UpdateMetadata(&h.metadata); err != nil {
			return err
		}
	}

	return nil
}

func (h *Heartbeat) updateMetaOnlineStatus(status bool) {
	h.shouldUpdateRepo = true
	h.metadata.Online = status
}

func (h *Heartbeat) updateLastMessageId(id int) {
	h.shouldUpdateRepo = true
	h.metadata.LastMessageId = id
}

func (h *Heartbeat) getLastGoodHeartbeatTimestamp() (time.Time, error) {
	resp, err := http.Get(h.config.HeartbeatUrl + "/heartbeats")
	if err != nil {
		return time.Time{}, nil
	}

	heartbeats, err := readHeartbeats(resp.Body)
	if err != nil {
		return time.Time{}, nil
	}

	lastGoodHeartbeat, err := h.getLastGoodHeartbeat(heartbeats)
	if err != nil {
		return time.Time{}, err
	}

	return lastGoodHeartbeat, nil
}

func (h *Heartbeat) getLastGoodHeartbeat(heartbeats []ReceivedHeartbeat) (time.Time, error) {
	for _, heartbeat := range heartbeats {
		if heartbeat.Name == h.config.HeartbeatName {
			return heartbeat.LastBeat, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot find good heartbeat")
}

func readHeartbeats(body io.ReadCloser) ([]ReceivedHeartbeat, error) {
	var hbr HeartBeatResponse
	err := json.NewDecoder(body).Decode(&hbr)

	return hbr.Heartbeats, err
}
