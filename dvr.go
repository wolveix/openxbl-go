package openxbl

import (
	"errors"
	"time"
)

const (
	DVRPrivacyBlocked        = DVRPrivacy("Blocked")
	DVRPrivacyEveryone       = DVRPrivacy("Everyone")
	DVRPrivacyPeopleOnMyList = DVRPrivacy("PeopleOnMyList")
)

type DVRPrivacy string

type GameClip struct {
	ContentID       string `json:"contentId"`
	ContentLocators []struct {
		Expiration  time.Time `json:"expiration,omitempty"`
		FileSize    int       `json:"fileSize,omitempty"`
		LocatorType string    `json:"locatorType"`
		Uri         string    `json:"uri"`
	} `json:"contentLocators"`
	ContentSegments []struct {
		SegmentID         int         `json:"segmentId"`
		CreationType      string      `json:"creationType"`
		CreatorChannelID  interface{} `json:"creatorChannelId"`
		CreatorXuid       int64       `json:"creatorXuid"`
		RecordDate        time.Time   `json:"recordDate"`
		DurationInSeconds int         `json:"durationInSeconds"`
		Offset            int         `json:"offset"`
		SecondaryTitleID  interface{} `json:"secondaryTitleId"`
		TitleID           int         `json:"titleId"`
	} `json:"contentSegments"`
	CreationType      string        `json:"creationType"`
	DurationInSeconds int           `json:"durationInSeconds"`
	FrameRate         int           `json:"frameRate"`
	GreatestMomentID  string        `json:"greatestMomentId"`
	LocalID           string        `json:"localId"`
	OwnerXuid         int64         `json:"ownerXuid"`
	ResolutionHeight  int           `json:"resolutionHeight"`
	ResolutionWidth   int           `json:"resolutionWidth"`
	SandboxID         string        `json:"sandboxId"`
	SharedTo          []interface{} `json:"sharedTo"`
	TitleData         string        `json:"titleData"`
	TitleID           int           `json:"titleId"`
	TitleName         string        `json:"titleName"`
	UploadDate        time.Time     `json:"uploadDate"`
	UploadLanguage    string        `json:"uploadLanguage"`
	UploadRegion      string        `json:"uploadRegion"`
	UploadTitleID     int           `json:"uploadTitleId"`
	UploadDeviceType  string        `json:"uploadDeviceType"`
	UserCaption       string        `json:"userCaption"`
	CommentCount      int           `json:"commentCount"`
	LikeCount         int           `json:"likeCount"`
	ShareCount        int           `json:"shareCount"`
	ViewCount         int           `json:"viewCount"`
	ContentState      string        `json:"contentState"`
	EnforcementState  string        `json:"enforcementState"`
	SafetyThreshold   string        `json:"safetyThreshold"`
	Sessions          []interface{} `json:"sessions"`
	Tournaments       []interface{} `json:"tournaments"`
}

func (c *Client) DeleteDVRGameClip(gameClipID string) error {
	if _, err := c.makeRequest("GET", "dvr/gameclips/delete/"+gameClipID, nil, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetDVRGameClips(continuationToken string) ([]*GameClip, string, error) {
	response := struct {
		ContinuationToken string      `json:"continuationToken"`
		GameClips         []*GameClip `json:"values"`
	}{}

	// if a continuation token (their version of pagination) is supplied, pass it to the API
	endpoint := "dvr/gameclips"
	if continuationToken != "" {
		endpoint += "?continuationToken=" + continuationToken
	}

	if _, err := c.makeRequest("GET", endpoint, nil, &response); err != nil {
		return nil, "", err
	}

	if len(response.GameClips) == 0 {
		return nil, "", errors.New("failed to find game clips")
	}

	return response.GameClips, response.ContinuationToken, nil
}

type Screenshot struct {
	CaptureDate     time.Time `json:"captureDate"`
	ContentID       string    `json:"contentId"`
	ContentLocators []struct {
		FileSize    int    `json:"fileSize,omitempty"`
		LocatorType string `json:"locatorType"`
		Uri         string `json:"uri"`
	} `json:"contentLocators"`
	CreationType     string        `json:"CreationType"`
	LocalID          string        `json:"localId"`
	OwnerXuid        int64         `json:"ownerXuid"`
	ResolutionHeight int           `json:"resolutionHeight"`
	ResolutionWidth  int           `json:"resolutionWidth"`
	SandboxID        string        `json:"sandboxId"`
	SharedTo         []interface{} `json:"sharedTo"`
	TitleID          int           `json:"titleId"`
	TitleName        string        `json:"titleName"`
	DateUploaded     time.Time     `json:"dateUploaded"`
	UploadLanguage   string        `json:"uploadLanguage"`
	UploadRegion     string        `json:"uploadRegion"`
	UploadTitleID    int           `json:"uploadTitleId"`
	UploadDeviceType string        `json:"uploadDeviceType"`
	CommentCount     int           `json:"commentCount"`
	LikeCount        int           `json:"likeCount"`
	ShareCount       int           `json:"shareCount"`
	ViewCount        int           `json:"viewCount"`
	ContentState     string        `json:"contentState"`
	EnforcementState string        `json:"enforcementState"`
	SafetyThreshold  string        `json:"safetyThreshold"`
	Sessions         []interface{} `json:"sessions"`
	Tournaments      []interface{} `json:"tournaments"`
}

func (c *Client) GetDVRScreenshots(continuationToken string) ([]*Screenshot, string, error) {
	response := struct {
		ContinuationToken string        `json:"continuationToken"`
		Screenshots       []*Screenshot `json:"values"`
	}{}

	// if a continuation token (their version of pagination) is supplied, pass it to the API
	endpoint := "dvr/screenshots"
	if continuationToken != "" {
		endpoint += "?continuationToken=" + continuationToken
	}

	if _, err := c.makeRequest("GET", endpoint, nil, &response); err != nil {
		return nil, "", err
	}

	if len(response.Screenshots) == 0 {
		return nil, "", errors.New("failed to find screenshots")
	}

	return response.Screenshots, response.ContinuationToken, nil
}

func (c *Client) SetDVRPrivacy(privacy DVRPrivacy) error {
	switch privacy {
	case DVRPrivacyBlocked, DVRPrivacyEveryone, DVRPrivacyPeopleOnMyList:
	default:
		return errors.New("invalid privacy type")
	}

	request := struct {
		Privacy string `json:"value"`
	}{
		Privacy: string(privacy),
	}

	if _, err := c.makeRequest("POST", "dvr/privacy", request, nil); err != nil {
		return err
	}

	return nil
}
