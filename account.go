package openxbl

import (
	"errors"
	"fmt"
	"strconv"
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

func (c *Client) GetAccount() (*Account, error) {
	response := struct {
		ProfileUsers []*Account `json:"profileUsers"`
	}{}

	if _, err := c.makeRequest("GET", "account", nil, &response); err != nil {
		return nil, err
	}

	if len(response.ProfileUsers) == 0 {
		return nil, errors.New("failed to find account info")
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
				return nil, fmt.Errorf("failed to convert gamerscore to int: %v", err)
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
