package heartbeat

import (
	"encoding/json"
	"fmt"
	"home_services_analyst/internal/config"
	"home_services_analyst/internal/repository"
	"home_services_analyst/internal/telegram"
	"home_services_analyst/internal/utils"

	"io"
	"net/http"
	"time"
)

const (
	GOOD_MESSAGE = "З'явилось світло 🎉"
	BAD_MESSAGE  = "Нема світла 🫠"

	STATUS_ONLINE  = true
	STATUS_OFFLINE = false
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

	shouldUpdateRepo bool
}

func NewHeartbeatService(
	cfg *config.HeartbeatConfig,
	repo *repository.Repository,
	tg *telegram.Telegram) Heartbeat {
	return Heartbeat{
		config:   cfg,
		repo:     repo,
		telegram: tg,
	}
}

func (h *Heartbeat) HandleHeartbeats() error {
	metadata, err := h.repo.LoadMetadata()
	if err != nil {
		return err
	}

	utils.PrettyPrint(metadata)

	lastGoodHeartbeat, err := h.getLastGoodHeartbeatTimestamp()
	if err != nil {
		return err
	}

	fmt.Println("Time since last heartbeat:", time.Since(lastGoodHeartbeat))

	lastHeartbeatInOkTime := time.Since(lastGoodHeartbeat) < time.Duration(h.config.CheckTimeoutSec*int(time.Second))

	fmt.Println("h.config.CheckTimeoutSec", time.Duration(h.config.CheckTimeoutSec))
	fmt.Println("lastHeartbeatInOkTime", lastHeartbeatInOkTime)
	fmt.Println("h.metadata.Online", metadata.Online)

	if lastHeartbeatInOkTime {
		if !metadata.Online {
			h.updateMetaOnlineStatus(&metadata, STATUS_ONLINE)

			fmt.Println("🚨🚨🚨 SEND GOOD MESSAGE 🚨🚨🚨")

			if _, err = h.telegram.SendTextMessage(GOOD_MESSAGE); err != nil {
				return err
			}
		}
	} else {
		if metadata.Online {
			h.updateMetaOnlineStatus(&metadata, STATUS_OFFLINE)

			fmt.Println("🚨🚨🚨 SEND BAD MESSAGE 🚨🚨🚨")

			resp, err := h.telegram.SendTextMessage(BAD_MESSAGE)
			if err != nil {
				return err
			}

			h.updateLastMessageId(&metadata, resp.Result.MessageId)
		} else {
			newMessage := fmt.Sprintf("%s\n\n%s", BAD_MESSAGE, fmt.Sprintf("%s%s", "Світла не було ", secondsToReadableUkrainian(int(time.Since(lastGoodHeartbeat).Seconds()))))
			h.telegram.UpdateTextMessage(metadata.LastBadMessageId, newMessage)
		}
	}

	if h.shouldUpdateRepo {
		// Use PostMetadata instead of actual metadata because of requirement from supabase - parameters
		// in post request should start with '_'
		if err := h.repo.UpdateMetadata(repository.PostMetadata(metadata)); err != nil {
			return err
		}
	}

	return nil
}

func (h *Heartbeat) updateMetaOnlineStatus(meta *repository.Metadata, status bool) {
	h.shouldUpdateRepo = true
	meta.Online = status
}

func (h *Heartbeat) updateLastMessageId(meta *repository.Metadata, id int) {
	h.shouldUpdateRepo = true
	meta.LastBadMessageId = id
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
