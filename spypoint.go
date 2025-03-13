package spypoint

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Client struct {
	Username string
	Password string
	UID      string
	Token    string
	client   *http.Client
}

func NewClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	if c.UID != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
}

func (c *Client) Login(force bool) error {
	if !force && c.UID != "" {
		return nil
	}
	payload := map[string]string{
		"username": c.Username,
		"password": c.Password,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://restapi.spypoint.com/api/v3/user/login", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("login failed")
	}

	var result struct {
		UID   string `json:"uuid"`
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	c.UID = result.UID
	c.Token = result.Token
	return nil
}
