package provider

import (
	"github.com/go-resty/resty/v2"
)

type CloudbackClient struct {
	restyClient *resty.Client
	Endpoint    string
	ApiKey      string
}

type BackupDefinition struct {
	Platform   string                   `json:"platform"`
	Account    string                   `json:"account"`
	Repository string                   `json:"repository"`
	Settings   BackupDefinitionSettings `json:"settings"`
}

type BackupDefinitionSettings struct {
	Enabled   bool   `json:"enabled"`
	Schedule  string `json:"schedule"`
	Storage   string `json:"storage"`
	Retention string `json:"retention"`
}

func NewCloudbackClient(baseURL, apiKey string) *CloudbackClient {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("X-API-KEY", apiKey)
	client.SetBaseURL(baseURL)
	client.SetDebug(false)

	return &CloudbackClient{
		restyClient: client,
		Endpoint:    baseURL,
		ApiKey:      apiKey,
	}
}

func (c *CloudbackClient) GetBackupDefinition(platform, account, repository string) (*BackupDefinition, error) {

	var response BackupDefinition

	resp, err := c.restyClient.R().
		SetBody(map[string]string{
			"platform":   platform,
			"account":    account,
			"repository": repository,
		}).
		SetResult(&response).
		Post("/api/ops/definition/get")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, NewAPIError(resp)
	}

	return &response, nil
}

func (c *CloudbackClient) UpdateBackupDefinition(platform, account, repository string, settings BackupDefinitionSettings) error {
	resp, err := c.restyClient.R().
		SetBody(&BackupDefinition{
			Platform:   platform,
			Account:    account,
			Repository: repository,
			Settings:   settings,
		}).
		Post("/api/ops/definition/update")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return NewAPIError(resp)
	}

	return nil
}

func NewAPIError(resp *resty.Response) error {
	return &APIError{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
	}
}

type APIError struct {
	StatusCode int
	Status     string
}

func (e *APIError) Error() string {
	return e.Status
}
