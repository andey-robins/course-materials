package shodan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type AlertHost struct {
	Name     string      `json:"name"`
	Filter   AlertFilter `json:"filters"`
	Lifetime int         `json:"expires"`
}

type AlertResponse struct {
	Name    string `json:"name"`
	Created string `json:"created"`
	// Triggers    []AlertTrigger `json:"triggers"`
	HasTriggers bool   `json:"has_triggers"`
	Expires     int    `json:"expires"`
	Expiration  string `json:"expiration"`
	// Filters     []AlertFilter  `json:"filters"`
	Id   string `json:"id"`
	Size int    `json:"size"`
}

type AddTriggerResponse struct {
	Success bool `json:"success"`
}

type AlertTrigger struct {
	Name        string `json:"name"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

type AlertFilter struct {
	IPs []string `json:"ip"`
}

// CreateAlert hits the /shodan/alert endpoint with a post request
// to create an alert on the host(s) defined by hostFilter.
func (s *Client) CreateAlert(name string, hostFilter AlertFilter, lifetime int) (*AlertResponse, error) {
	alertBody := AlertHost{
		Name:     name,
		Filter:   hostFilter,
		Lifetime: lifetime,
	}

	alertBodyJson, err := json.Marshal(alertBody)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(fmt.Sprintf("%s/shodan/alert?key=%s", BaseURL, s.apiKey), "application/json", bytes.NewReader(alertBodyJson))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ret AlertResponse
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

// GetMyAlerts returns a list of all alerts created by the user
// by hitting /shodan/alert/info
func (s *Client) GetMyAlerts() ([]*AlertResponse, error) {
	ret := make([]*AlertResponse, 0)

	res, err := http.Get(fmt.Sprintf("%s/shodan/alert/info?key=%s", BaseURL, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// GetAllValidTriggers returns a list of all valid triggers as
// defined by the /shodan/alert/triggers endpoint
func (s *Client) GetAllValidTriggers() ([]*AlertTrigger, error) {
	// TODO: Cache these for a period of time?
	ret := make([]*AlertTrigger, 0)

	res, err := http.Get(fmt.Sprintf("%s/shodan/alert/triggers?key=%s", BaseURL, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// AddTriggerToAlert adds a trigger to an alert by hitting the
// /shodan/alert/{id}/trigger/{trigger} endpoint
func (s *Client) AddTriggerToAlert(triggerName, alertId string) error {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/shodan/alert/%s/trigger/%s?key=%s", BaseURL, alertId, triggerName, s.apiKey), nil)
	if err != nil {
		log.Panicln(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer res.Body.Close()

	var ret AddTriggerResponse
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		log.Panicln(err)
	}

	if !ret.Success {
		return errors.New("failed to add trigger to alert")
	}

	return nil
}
