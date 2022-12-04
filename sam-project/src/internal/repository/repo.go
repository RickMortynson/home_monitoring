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
	supaUrl    string
	supaBucket string
	supaToken  string
	filename   string
}

type Metadata struct {
	Online           bool   `json:"online"`
	ShutdownTableMd5 string `json:"shutdownTableMd5"`
	LastMessageId    int    `json:"lastMessageId"`
}

func NewRepository(repo *config.RepoConfig) Repository {
	return Repository{
		supaUrl:    repo.SupabaseUrl,
		supaBucket: repo.SupabaseBucket,
		supaToken:  repo.SupabaseToken,
		filename:   repo.MetadataFile,
	}
}

func (r *Repository) LoadMetadata(meta *Metadata) error {
	url := fmt.Sprintf("%s/storage/v1/object/authenticated/%s/%s", r.supaUrl, r.supaBucket, r.filename)

	fmt.Println("Send url", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+r.supaToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	responseBody, _ := io.ReadAll(res.Body)

	if err = json.Unmarshal(responseBody, &meta); err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateMetadata(meta *Metadata) error {
	fmt.Println("Updating metadata...", meta)
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", r.supaUrl, r.supaBucket, r.filename)
	
	fmt.Println("Update url", url)

	reqBody, err := json.Marshal(*meta)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+r.supaToken)

	if _, err = http.DefaultClient.Do(req); err != nil {
		return err
	}

	return nil
}
