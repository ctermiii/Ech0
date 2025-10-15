package fediverse

import (
	"encoding/json"
	"errors"
	"net/http"

	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

//==============================================================================
//	Fetch
//==============================================================================

// FetchRemoteActorInbox 获取远程 Actor 的 Inbox URL
func (core *FediverseCore) FetchRemoteActorInbox(actorURL string) (string, error) {
	if actorURL == "" {
		return "", errors.New("remote actor url is empty")
	}

	body, err := httpUtil.SendRequest(actorURL, http.MethodGet, httpUtil.Header{
		Header:  "Accept",
		Content: "application/activity+json",
	})
	if err != nil {
		return "", err
	}

	var resp struct {
		Inbox     string `json:"inbox"`
		Endpoints struct {
			SharedInbox string `json:"sharedInbox"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	if resp.Inbox != "" {
		return resp.Inbox, nil
	}
	if resp.Endpoints.SharedInbox != "" {
		return resp.Endpoints.SharedInbox, nil
	}

	return "", errors.New("remote actor inbox not found")
}
