package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"home_services_analyst/internal/config"
	"io"
	"net/http"
)

type Repository struct {
	supaUrl   string
	supaToken string
}

type Metadata struct {
	Online           bool   `json:"online"`
	PowerScheduleMd5 string `json:"power_schedule_md5"`
	LastBadMessageId int    `json:"last_bad_message_id"`
}

type PostMetadata struct {
	Online           bool   `json:"_online"`
	PowerScheduleMd5 string `json:"_power_schedule_md5"`
	LastBadMessageId int    `json:"_last_bad_message_id"`
}

func NewRepository(repo *config.RepoConfig) Repository {
	return Repository{
		supaUrl:   repo.SupabaseUrl,
		supaToken: repo.SupabaseToken,
	}
}

func (r *Repository) LoadMetadata() (Metadata, error) {
	url := fmt.Sprintf("%s/rest/v1/rpc/get_metadata", r.supaUrl)

	fmt.Println("Send url", url)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return Metadata{}, err
	}

	req.Header.Add("apikey", r.supaToken)
	req.Header.Add("Authorization", "Bearer "+r.supaToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Metadata{}, err
	}

	responseBody, _ := io.ReadAll(res.Body)

	fmt.Println("RESPONSE🔔", string(responseBody))

	var metaResponse []Metadata
	if err = json.Unmarshal(responseBody, &metaResponse); err != nil {
		return Metadata{}, err
	}

	fmt.Println("UNMARSHALLED RESPONSE🔔", metaResponse)

	return metaResponse[0], nil
}

func (r *Repository) UpdateMetadata(meta PostMetadata) error {
	fmt.Println("Updating metadata...", meta)
	url := fmt.Sprintf("%s/rest/v1/rpc/update_metadata", r.supaUrl)

	reqBody, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("apikey", r.supaToken)
	req.Header.Add("Authorization", "Bearer "+r.supaToken)

	if _, err = http.DefaultClient.Do(req); err != nil {
		return err
	}

	return nil
}
