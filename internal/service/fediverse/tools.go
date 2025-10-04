package service

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

	"github.com/lin-snow/ech0/internal/config"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	settingModel "github.com/lin-snow/ech0/internal/model/setting"
	userModel "github.com/lin-snow/ech0/internal/model/user"
	fileUtil "github.com/lin-snow/ech0/internal/util/file"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

//==============================================================================
//	normalize or resolve or generate
//==============================================================================

// normalizeServerURL 标准化服务器 URL，确保有协议头且无尾部斜杠
func normalizeServerURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New(commonModel.ACTIVEPUB_NOT_ENABLED)
	}
	if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
		trimmed = "https://" + trimmed
	}
	return strings.TrimRight(trimmed, "/"), nil
}

// normalizePageParams 标准化分页参数
func normalizePageParams(page, pageSize int) (int, int) {
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

// resolveActorURL 解析输入，返回 Actor URL，格式为 http(s)://domain/users/username
func resolveActorURL(input string) (string, error) {
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
	webfingerURL := fmt.Sprintf("https://%s/.well-known/webfinger?resource=%s", domain, url.QueryEscape("acct:"+username+"@"+domain))
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
			if link.Type == "application/activity+json" || link.Type == "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" || link.Type == "" {
				// 返回找到的 Actor URL,格式为 http(s)://domain/users/username
				return link.Href, nil
			}
		}
	}

	return "", errors.New(commonModel.GET_ACTOR_ERROR)
}

// generateDeterministicActivityID 生成确定性的 Activity ID
func generateDeterministicActivityID(serverURL, username, prefix, key string) string {
	hash := sha256.Sum256([]byte(strings.ToLower(key)))
	short := hex.EncodeToString(hash[:16])
	return fmt.Sprintf("%s/activities/%s/%s/%s", serverURL, username, prefix, short)
}

//==============================================================================
//	Convert
//==============================================================================

// ConvertEchoToActivity 将 Echo 转换为 ActivityPub Activity
func (fediverseService *FediverseService) ConvertEchoToActivity(echo *echoModel.Echo, actor *model.Actor, serverURL string) model.Activity {
	obj := fediverseService.ConvertEchoToObject(echo, actor, serverURL)

	activityID := fmt.Sprintf("%s/activities/%d", serverURL, echo.ID)

	activity := model.Activity{
		Context: []any{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		ActivityID: activityID,
		Type:       model.ActivityTypeCreate,
		ActorID:    actor.ID,
		ActorURL:   actor.ID,
		ObjectID:   obj.ObjectID,
		ObjectType: obj.Type,
		Published:  echo.CreatedAt,
		To:         obj.To,
		Cc:         []string{actor.Followers},
		Summary:    "",
		Delivered:  false,
		CreatedAt:  time.Now(),
	}

	activityJSON, _ := json.Marshal(activity)
	activity.ActivityJSON = string(activityJSON)
	return activity
}

// ConvertEchoToObject 将 Echo 转换为 ActivityPub Object
func (fediverseService *FediverseService) ConvertEchoToObject(echo *echoModel.Echo, actor *model.Actor, serverURL string) model.Object {
	var attachments []model.Attachment
	for i := range echo.Images {
		attachments = append(attachments, model.Attachment{
			Type:      "Image",
			MediaType: httpUtil.GetMIMETypeFromFilenameOrURL(echo.Images[i].ImageURL),
			URL:       fileUtil.GetImageURL(echo.Images[i], serverURL),
		})
	}

	return model.Object{
		Context: []any{
			"https://www.w3.org/ns/activitystreams",
		},
		ObjectID:     fmt.Sprintf("%s/objects/%d", serverURL, echo.ID),
		Type:         "Note",
		Content:      echo.Content,
		AttributedTo: actor.ID,
		Published:    echo.CreatedAt,
		To: []string{
			"https://www.w3.org/ns/activitystreams#Public",
		},
		Attachments: attachments,
	}
}

//==============================================================================
// Build
//==============================================================================

// BuildActor 构建 Actor 对象
func (fediverseService *FediverseService) BuildActor(user *userModel.User) (model.Actor, *settingModel.SystemSetting, error) {
	// 从设置服务获取服务器域名
	var setting settingModel.SystemSetting
	if err := fediverseService.settingService.GetSetting(&setting); err != nil {
		return model.Actor{}, nil, err
	}
	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return model.Actor{}, nil, err
	}
	// 构建头像信息 (域名 + /api + 头像路径)
	if user.Avatar == "" {
		user.Avatar = "/Ech0.png" // 默认头像路径
	} else {
		user.Avatar = "/api" + user.Avatar
	}
	avatarURL := serverURL + user.Avatar
	avatarMIME := httpUtil.GetMIMETypeFromFilenameOrURL(avatarURL)

	// 构建 Actor 对象
	return model.Actor{
		Context: []any{
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
		},
		ID:                serverURL + "/users/" + user.Username, // 实例地址拼接 域名 + /users/ + username
		Type:              "Person",                              // 固定值
		Name:              setting.ServerName,                    // 显示名称
		PreferredUsername: user.Username,                         // 用户名
		Summary:           "你好呀!👋 我是来自Ech0的" + user.Username,     // 简介
		Icon: model.Preview{
			Type:      "Image",
			MediaType: avatarMIME,
			URL:       avatarURL,
		},
		Image: model.Preview{
			Type:      "Image",
			MediaType: "image/png",
			URL:       serverURL + "/banner.png", // 封面图片，固定为 /banner.png
		},
		Followers: serverURL + "/users/" + user.Username + "/followers", // 粉丝列表地址
		Following: serverURL + "/users/" + user.Username + "/following", // 关注列表地址
		Inbox:     serverURL + "/users/" + user.Username + "/inbox",     // 收件箱地址
		Outbox:    serverURL + "/users/" + user.Username + "/outbox",    // 发件箱地址
		PublicKey: model.PublicKey{
			ID:           serverURL + "/users/" + user.Username + "#main-key",
			Owner:        serverURL + "/users/" + user.Username,
			PublicKeyPem: string(config.RSA_PUBLIC_KEY),
			Type:         "Key",
		},
	}, &setting, nil
}

// BuildOutbox 构建 Outbox 元信息
func (fediverseService *FediverseService) BuildOutbox(username string) (model.OutboxResponse, error) {
	// 查询用户，确保用户存在
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.OutboxResponse{}, errors.New(commonModel.USER_NOTFOUND)
	}

	// 获取 Actor和 setting
	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.OutboxResponse{}, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return model.OutboxResponse{}, err
	}

	// 查 Echos
	_, total := fediverseService.echoRepository.GetEchosByPage(1, 10, "", false)

	firstPage := fmt.Sprintf("%s?page=1", actor.Outbox)
	lastPage := ""
	if total > 0 {
		totalPages := int(total) / 10
		if total%10 != 0 {
			totalPages++
		}
		lastPage = fmt.Sprintf("%s?page=%d", actor.Outbox, totalPages)
	}

	return model.OutboxResponse{
		Context:    "https://www.w3.org/ns/activitystreams",
		ID:         fmt.Sprintf("%s/users/%s/outbox", serverURL, username),
		Type:       "OrderedCollection",
		TotalItems: int(total),
		First:      firstPage,
		Last:       lastPage,
	}, nil
}

// buildAcceptActivityPayload 构建 Accept Activity 的 JSON Payload
func (fediverseService *FediverseService) buildAcceptActivityPayload(actor *model.Actor, follow *model.Activity, followerActor, serverURL string) ([]byte, error) {
	if follow.ActivityID == "" {
		return nil, errors.New("follow activity missing id")
	}

	target := follow.ObjectID
	if target == "" {
		target = actor.ID
	}

	now := time.Now().UTC()
	acceptID := fmt.Sprintf("%s/activities/%s/accept/%d", serverURL, actor.PreferredUsername, now.UnixNano())

	payload := map[string]any{
		"@context": []any{"https://www.w3.org/ns/activitystreams"},
		"id":       acceptID,
		"type":     model.ActivityTypeAccept,
		"actor":    actor.ID,
		"object": map[string]any{
			"id":     follow.ActivityID,
			"type":   model.ActivityTypeFollow,
			"actor":  followerActor,
			"object": target,
		},
		"to":        []string{followerActor},
		"published": now.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

// buildFollowActivityPayload 构建 Follow Activity 的 JSON Payload
func buildFollowActivityPayload(actor *model.Actor, targetActor string, activityID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if activityID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" {
		return nil, errors.New("target actor is empty")
	}

	payload := map[string]any{
		"@context":  []any{"https://www.w3.org/ns/activitystreams"},
		"id":        activityID,
		"type":      model.ActivityTypeFollow,
		"actor":     actor.ID,
		"object":    targetActor,
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

// buildUndoFollowActivityPayload 构建 Undo Follow Activity 的 JSON Payload
func buildUndoFollowActivityPayload(actor *model.Actor, targetActor string, undoID string, followActivityID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if undoID == "" || followActivityID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" {
		return nil, errors.New("target actor is empty")
	}

	payload := map[string]any{
		"@context": []any{"https://www.w3.org/ns/activitystreams"},
		"id":       undoID,
		"type":     model.ActivityTypeUndo,
		"actor":    actor.ID,
		"object": map[string]any{
			"id":     followActivityID,
			"type":   model.ActivityTypeFollow,
			"actor":  actor.ID,
			"object": targetActor,
		},
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

// buildLikeActivityPayload 构建 Like Activity 的 JSON Payload
func buildLikeActivityPayload(actor *model.Actor, targetActor string, object string, activityID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if activityID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" || object == "" {
		return nil, errors.New("target actor or object is empty")
	}

	payload := map[string]any{
		"@context":  []any{"https://www.w3.org/ns/activitystreams"},
		"id":        activityID,
		"type":      model.ActivityTypeLike,
		"actor":     actor.ID,
		"object":    object,
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

// buildUndoLikeActivityPayload 构建 Undo Like Activity 的 JSON Payload
func buildUndoLikeActivityPayload(actor *model.Actor, targetActor string, object string, likeActivityID string, undoID string, published time.Time) ([]byte, error) {
	if actor == nil {
		return nil, errors.New("actor is nil")
	}
	if likeActivityID == "" || undoID == "" {
		return nil, errors.New("activity id is empty")
	}
	if targetActor == "" || object == "" {
		return nil, errors.New("target actor or object is empty")
	}

	payload := map[string]any{
		"@context": []any{"https://www.w3.org/ns/activitystreams"},
		"id":       undoID,
		"type":     model.ActivityTypeUndo,
		"actor":    actor.ID,
		"object": map[string]any{
			"id":     likeActivityID,
			"type":   model.ActivityTypeLike,
			"actor":  actor.ID,
			"object": object,
		},
		"to":        []string{targetActor},
		"published": published.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

//==============================================================================
//	Fetch or Push
//==============================================================================

// fetchRemoteActorInbox 获取远程 Actor 的 Inbox URL
func (fediverseService *FediverseService) fetchRemoteActorInbox(actorURL string) (string, error) {
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

// PushEchoToFediverse 将 Echo 推送到联邦网络
func (fediverseService *FediverseService) PushEchoToFediverse(userId uint, echo echoModel.Echo) error {
	// 获取用户
	user, err := fediverseService.commonService.CommonGetUserByUserId(userId)
	if err != nil {
		return err
	}

	// 获取粉丝列表
	followers, err := fediverseService.fediverseRepository.GetFollowers(user.ID)
	if err != nil {
		return err
	}
	if len(followers) == 0 {
		return nil
	}

	// 获取 Actor 和 setting
	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return err
	}

	activity := fediverseService.ConvertEchoToActivity(&echo, &actor, serverURL)
	object := fediverseService.ConvertEchoToObject(&echo, &actor, serverURL)

	activityMap := map[string]any{}
	activityBytes, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(activityBytes, &activityMap); err != nil {
		return err
	}

	objectMap := map[string]any{}
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(objectBytes, &objectMap); err != nil {
		return err
	}

	activityMap["object"] = objectMap

	payloadBytes, err := json.Marshal(activityMap)
	if err != nil {
		return err
	}

	var errs []error
	// 推送到每个粉丝的Inbox
	for _, follower := range followers {
		inboxURL, err := fediverseService.fetchRemoteActorInbox(follower.ActorID)
		if err != nil {
			errs = append(errs, fmt.Errorf("fetch inbox for %s: %w", follower.ActorID, err))
			continue
		}

		if err := httpUtil.PostActivity(payloadBytes, inboxURL, actor.ID); err != nil {
			errs = append(errs, fmt.Errorf("post activity to %s: %w", inboxURL, err))
		}
	}

	if len(errs) > 0 {
		fmt.Println("Errors occurred while pushing to Fediverse:")
		for _, e := range errs {
			fmt.Println(e)
		}
		return errors.Join(errs...)
	}

	return nil
}
