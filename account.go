package openxbl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	AccountTierGold   = "Gold"
	AccountTierSilver = "Silver"
)

type Account struct {
	ID                string `json:"id" yaml:"id"`
	HostID            string `json:"hostId" yaml:"hostId"`
	AvatarURL         string `json:"avatarUrl" yaml:"avatarUrl"`
	Bio               string `json:"bio" yaml:"bio"`
	Gamerscore        int    `json:"gamerscore" yaml:"gamerscore"`
	Gamertag          string `json:"gamertag" yaml:"gamertag"`
	IsSponsoredUser   bool   `json:"isSponsoredUser"`
	Location          string `json:"location" yaml:"location"`
	PreferredColorURL string `json:"preferredColorURL" yaml:"preferredColorURL"`
	RealName          string `json:"realName" yaml:"realName"`
	Tier              string `json:"tier" yaml:"tier"`
	RawSettings       []struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"settings"`
}

func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	response := struct {
		ProfileUsers []*Account `json:"profileUsers"`
	}{}

	if _, err := c.makeRequest(ctx, http.MethodGet, "account", nil, &response); err != nil {
		return nil, err
	}

	if len(response.ProfileUsers) == 0 {
		return nil, errors.New("find account info")
	}

	account := response.ProfileUsers[0]

	for _, setting := range account.RawSettings {
		switch setting.ID {
		case "AccountTier":
			account.Tier = setting.Value
		case "Bio":
			account.Bio = setting.Value
		case "GameDisplayPicRaw":
			account.AvatarURL = setting.Value
		case "Gamerscore":
			gamerscore, err := strconv.Atoi(setting.Value)
			if err != nil {
				return nil, fmt.Errorf("convert gamerscore to int: %w", err)
			}

			account.Gamerscore = gamerscore
		case "Gamertag":
			account.Gamertag = setting.Value
		case "Location":
			account.Location = setting.Value
		case "PreferredColor":
			account.PreferredColorURL = setting.Value
		case "RealName":
			account.RealName = setting.Value
		}
	}

	return response.ProfileUsers[0], nil
}

// GenerateGamertags returns a list of generated gamertag options.
func (c *Client) GenerateGamertags(ctx context.Context, quantity int) ([]string, error) {
	if quantity <= 0 {
		return nil, errors.New("invalid quantity")
	}

	request := struct {
		Algorithm int    `json:"algorithm"`
		Count     int    `json:"count"`
		Seed      string `json:"seed"`
		Locale    string `json:"locale"`
	}{1, quantity, "", "en-US"}

	response := struct {
		Gamertags []string `json:"Gamertags"`
	}{}

	if _, err := c.makeRequest(ctx, http.MethodPost, "generate/gamertag", request, &response); err != nil {
		return nil, err
	}

	if len(response.Gamertags) == 0 {
		return nil, errors.New("no gamertags generated")
	}

	return response.Gamertags, nil
}

type Presence struct {
	ID      string `json:"xuid"`
	Devices []struct {
		Type   string `json:"type"`
		Titles []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Placement    string `json:"placement"`
			State        string `json:"state"`
			LastModified string `json:"lastModified"`
		} `json:"titles"`
	} `json:"devices"`
	LastSeen struct {
		DeviceType string `json:"deviceType"`
		TitleID    string `json:"titleId"`
		TitleName  string `json:"titleName"`
		Timestamp  string `json:"timestamp"`
	} `json:"lastSeen"`
	State string `json:"state"`
}

// GetPresenceForUser returns the current Presence for the given user ID.
func (c *Client) GetPresenceForUser(ctx context.Context, xboxIDs ...string) ([]*Presence, error) {
	if len(xboxIDs) == 0 {
		return nil, errors.New("missing xbox ID")
	}

	var response []*Presence

	if _, err := c.makeRequest(ctx, http.MethodGet, strings.Join(xboxIDs, ",")+"/presence", nil, &response); err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, errors.New("find user presences")
	}

	return response, nil
}
