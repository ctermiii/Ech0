package fediverse

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

//==============================================================================
//	normalize or resolve or generate
//==============================================================================

// NormalizeServerURL 标准化服务器 URL，确保有协议头且无尾部斜杠
func NormalizeServerURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New(commonModel.ACTIVEPUB_NOT_ENABLED)
	}
	if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
		trimmed = "https://" + trimmed
	}
	return strings.TrimRight(trimmed, "/"), nil
}

// NormalizePageParams 标准化分页参数
func NormalizePageParams(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = model.DefaultCollectionPageSize
	} else if pageSize > model.MaxCollectionPageSize {
		pageSize = model.MaxCollectionPageSize
	}
	return page, pageSize
}

// ResolveActorURL 解析输入，返回 Actor URL，格式为 http(s)://domain/users/username
func ResolveActorURL(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed, nil
	}

	// 处理 acct:username@domain 或 username@domain 格式
	resource := trimmed
	if after, ok := strings.CutPrefix(resource, "acct:"); ok {
		resource = after
	}
	resource = strings.TrimPrefix(resource, "@")

	// 必须包含 '@' 分隔符
	if !strings.Contains(resource, "@") {
		return "", errors.New(commonModel.GET_ACTOR_ERROR)
	}

	// 分割用户名和域名
	parts := strings.SplitN(resource, "@", 2)
	username := strings.TrimSpace(parts[0])
	domain := strings.TrimSpace(parts[1])
	if username == "" || domain == "" {
		return "", errors.New(commonModel.GET_ACTOR_ERROR)
	}

	// 通过 WebFinger 获取 Actor URL
	webfingerURL := fmt.Sprintf(
		"https://%s/.well-known/webfinger?resource=%s",
		domain,
		url.QueryEscape("acct:"+username+"@"+domain),
	)
	body, err := httpUtil.SendRequest(webfingerURL, http.MethodGet, httpUtil.Header{
		Header:  "Accept",
		Content: "application/jrd+json, application/json",
	}, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("%s: %w", commonModel.GET_ACTOR_ERROR, err)
	}

	var resp struct {
		Links []struct {
			Rel  string `json:"rel"`
			Type string `json:"type"`
			Href string `json:"href"`
		} `json:"links"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("%s: %w", commonModel.GET_ACTOR_ERROR, err)
	}

	// 查找符合条件的 self 链接
	for _, link := range resp.Links {
		if link.Rel == "self" && link.Href != "" {
			if link.Type == "application/activity+json" ||
				link.Type == "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" ||
				link.Type == "" {
				// 返回找到的 Actor URL,格式为 http(s)://domain/users/username
				return link.Href, nil
			}
		}
	}

	return "", errors.New(commonModel.GET_ACTOR_ERROR)
}

// GenerateDeterministicActivityID 生成确定性的 Activity ID
func GenerateDeterministicActivityID(serverURL, username, prefix, key string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(key)))
	short := hex.EncodeToString(hash[:16])
	return fmt.Sprintf("%s/activities/%s/%s/%s", serverURL, username, prefix, short)
}
