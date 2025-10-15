package service

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"

// 	model "github.com/lin-snow/ech0/internal/model/fediverse"
// 	userModel "github.com/lin-snow/ech0/internal/model/user"
// 	httpUtil "github.com/lin-snow/ech0/internal/util/http"
// 	mdUtil "github.com/lin-snow/ech0/internal/util/md"
// )

// const inboxFetchTimeout = 5 * time.Second

// // handleCreateActivity 处理接收到的 Create 活动，将远端推文落库以便时间线展示
// func (fediverse *FediverseService) handleCreateActivity(user *userModel.User, activity *model.Activity) error {
// 	if user == nil {
// 		return errors.New("user is nil")
// 	}
// 	if activity == nil {
// 		return errors.New("activity is nil")
// 	}

// 	remoteActor := strings.TrimSpace(activity.ActorURL)
// 	if remoteActor == "" {
// 		remoteActor = strings.TrimSpace(activity.ActorID)
// 	}
// 	if remoteActor == "" {
// 		return errors.New("create activity missing actor")
// 	}

// 	objectMap, objectJSON, objectID, objectType, attributedTo, err := fediverse.resolveActivityObject(activity)
// 	if err != nil {
// 		return err
// 	}

// 	activityID := strings.TrimSpace(activity.ActivityID)
// 	if activityID == "" {
// 		activityID = objectID
// 	}
// 	if activityID == "" {
// 		return errors.New("create activity missing id")
// 	}

// 	content := normalizeActivityContent(getStringFromMap(objectMap, "content"), objectMap)
// 	summary := strings.TrimSpace(activity.Summary)
// 	if summary == "" {
// 		summary = getStringFromMap(objectMap, "summary")
// 	}
// 	actorDisplayName := getStringFromMap(objectMap, "name")
// 	preferredUsername := derivePreferredUsername(remoteActor)
// 	if actorDisplayName == "" {
// 		actorDisplayName = preferredUsername
// 	}

// 	avatarURL := fediverse.resolveActorAvatar(activity, objectMap)

// 	publishedAt := fediverse.resolvePublishedAt(activity, objectMap)
// 	toJSON := mustMarshalStrings(activity.To)
// 	ccJSON := mustMarshalStrings(activity.Cc)

// 	activity.ObjectID = objectID
// 	activity.ObjectType = objectType

// 	activityJSONBytes, err := json.Marshal(activity)
// 	if err != nil {
// 		return fmt.Errorf("marshal activity: %w", err)
// 	}
// 	activity.ActivityJSON = string(activityJSONBytes)

// 	status := &model.InboxStatus{
// 		UserID:                 user.ID,
// 		ActivityID:             activityID,
// 		ActorID:                remoteActor,
// 		ActorPreferredUsername: preferredUsername,
// 		ActorDisplayName:       actorDisplayName,
// 		ActorAvatar:            avatarURL,
// 		ObjectID:               objectID,
// 		ObjectType:             objectType,
// 		ObjectAttributedTo:     attributedTo,
// 		Summary:                summary,
// 		Content:                content,
// 		To:                     toJSON,
// 		Cc:                     ccJSON,
// 		RawActivity:            string(activityJSONBytes),
// 		RawObject:              string(objectJSON),
// 		PublishedAt:            publishedAt,
// 	}

// 	return fediverse.txManager.Run(func(ctx context.Context) error {
// 		return fediverse.fediverseRepository.UpsertInboxStatus(ctx, status)
// 	})
// }

// // resolveActivityObject 解析 Activity 中的 Object 字段，必要时远程抓取完整的 Object
// func (fediverse *FediverseService) resolveActivityObject(
// 	activity *model.Activity,
// ) (map[string]any, []byte, string, string, string, error) {
// 	var (
// 		objectMap map[string]any
// 		objectID  string
// 		objectTyp string
// 		attTo     string
// 	)

// 	switch obj := activity.Object.(type) {
// 	case string:
// 		objectID = strings.TrimSpace(obj)
// 	case map[string]any:
// 		objectMap = obj
// 	case []any:
// 		if len(obj) > 0 {
// 			if candidate, ok := obj[0].(map[string]any); ok {
// 				objectMap = candidate
// 			}
// 		}
// 	}

// 	if objectMap != nil {
// 		objectID = strings.TrimSpace(getStringFromMap(objectMap, "id"))
// 		if objectID == "" {
// 			objectID = strings.TrimSpace(getStringFromMap(objectMap, "url"))
// 		}
// 		objectTyp = getStringFromMap(objectMap, "type")
// 		attTo = extractAttributedTo(objectMap["attributedTo"])
// 	}

// 	var objectJSON []byte
// 	if objectMap == nil && objectID != "" {
// 		body, err := httpUtil.SendRequest(objectID, http.MethodGet, httpUtil.Header{
// 			Header:  "Accept",
// 			Content: "application/activity+json",
// 		}, inboxFetchTimeout)
// 		if err != nil {
// 			return nil, nil, "", "", "", fmt.Errorf("fetch object %s: %w", objectID, err)
// 		}
// 		if err := json.Unmarshal(body, &objectMap); err != nil {
// 			return nil, nil, "", "", "", fmt.Errorf("decode object %s: %w", objectID, err)
// 		}
// 		objectJSON = body
// 		if objectTyp == "" {
// 			objectTyp = getStringFromMap(objectMap, "type")
// 		}
// 		if attTo == "" {
// 			attTo = extractAttributedTo(objectMap["attributedTo"])
// 		}
// 	}

// 	if objectMap == nil && objectID == "" {
// 		return nil, nil, "", "", "", errors.New("create activity missing object id")
// 	}
// 	if objectMap == nil {
// 		objectMap = map[string]any{"id": objectID}
// 	}

// 	if objectID == "" {
// 		objectID = strings.TrimSpace(getStringFromMap(objectMap, "id"))
// 	}
// 	if objectID == "" {
// 		return nil, nil, "", "", "", errors.New("create activity missing object id")
// 	}

// 	if objectTyp == "" {
// 		objectTyp = getStringFromMap(objectMap, "type")
// 	}

// 	if len(objectJSON) == 0 {
// 		serialized, err := json.Marshal(objectMap)
// 		if err != nil {
// 			return nil, nil, "", "", "", fmt.Errorf("marshal object %s: %w", objectID, err)
// 		}
// 		objectJSON = serialized
// 	}

// 	return objectMap, objectJSON, objectID, objectTyp, attTo, nil
// }

// // resolvePublishedAt 从 Activity 和 Object 中推断发布时间
// func (fediverse *FediverseService) resolvePublishedAt(activity *model.Activity, objectMap map[string]any) time.Time {
// 	if activity.Published != (time.Time{}) {
// 		return activity.Published
// 	}

// 	if candidate := getStringFromMap(objectMap, "published"); candidate != "" {
// 		if ts, err := parseRFC3339(candidate); err == nil {
// 			return ts
// 		}
// 	}

// 	return time.Now().UTC()
// }

// // getStringFromMap 安全地从 map 中读取字符串字段
// func getStringFromMap(payload map[string]any, key string) string {
// 	if payload == nil {
// 		return ""
// 	}
// 	if value, ok := payload[key]; ok {
// 		return extractString(value)
// 	}
// 	return ""
// }

// // extractAttributedTo 处理 attributedTo 字段，可能是字符串、对象或数组
// func extractAttributedTo(value any) string {
// 	switch v := value.(type) {
// 	case string:
// 		return strings.TrimSpace(v)
// 	case map[string]any:
// 		if id := getStringFromMap(v, "id"); id != "" {
// 			return id
// 		}
// 		return getStringFromMap(v, "url")
// 	case []any:
// 		for _, item := range v {
// 			if s := extractAttributedTo(item); s != "" {
// 				return s
// 			}
// 		}
// 	}
// 	return ""
// }

// func (fediverse *FediverseService) resolveActorAvatar(activity *model.Activity, objectMap map[string]any) string {
// 	if icon := extractIconURL(objectMap, "icon"); icon != "" {
// 		return icon
// 	}
// 	if icon := extractIconURL(objectMap, "image"); icon != "" {
// 		return icon
// 	}
// 	if icon := extractIconURL(map[string]any{"value": objectMap["attributedTo"]}, "value"); icon != "" {
// 		return icon
// 	}
// 	if icon := extractIconURL(map[string]any{"value": objectMap["actor"]}, "value"); icon != "" {
// 		return icon
// 	}

// 	checked := make(map[string]struct{})
// 	tryFetch := func(candidate string) string {
// 		candidate = strings.TrimSpace(candidate)
// 		if candidate == "" {
// 			return ""
// 		}
// 		if _, exists := checked[candidate]; exists {
// 			return ""
// 		}
// 		checked[candidate] = struct{}{}

// 		icon, err := fediverse.fetchActorIcon(candidate)
// 		if err != nil {
// 			return ""
// 		}
// 		return icon
// 	}

// 	for _, candidate := range extractActorCandidates(objectMap["attributedTo"]) {
// 		if icon := tryFetch(candidate); icon != "" {
// 			return icon
// 		}
// 	}

// 	for _, candidate := range extractActorCandidates(objectMap["actor"]) {
// 		if icon := tryFetch(candidate); icon != "" {
// 			return icon
// 		}
// 	}

// 	remoteActor := strings.TrimSpace(activity.ActorURL)
// 	if remoteActor == "" {
// 		remoteActor = strings.TrimSpace(activity.ActorID)
// 	}
// 	if icon := tryFetch(remoteActor); icon != "" {
// 		return icon
// 	}

// 	return ""
// }

// func (fediverse *FediverseService) fetchActorIcon(actorURL string) (string, error) {
// 	actorURL = strings.TrimSpace(actorURL)
// 	if actorURL == "" {
// 		return "", errors.New("actor url is empty")
// 	}

// 	body, err := httpUtil.SendRequest(actorURL, http.MethodGet, httpUtil.Header{
// 		Header:  "Accept",
// 		Content: "application/activity+json",
// 	}, inboxFetchTimeout)
// 	if err != nil {
// 		return "", err
// 	}

// 	var actor map[string]any
// 	if err := json.Unmarshal(body, &actor); err != nil {
// 		return "", err
// 	}

// 	if icon := extractIconURL(actor, "icon"); icon != "" {
// 		return icon, nil
// 	}
// 	if icon := extractIconURL(actor, "image"); icon != "" {
// 		return icon, nil
// 	}

// 	return "", nil
// }

// func extractIconURL(container map[string]any, key string) string {
// 	if container == nil {
// 		return ""
// 	}
// 	value, ok := container[key]
// 	if !ok {
// 		return ""
// 	}
// 	return extractIconValue(value)
// }

// func extractIconValue(value any) string {
// 	switch v := value.(type) {
// 	case string:
// 		candidate := strings.TrimSpace(v)
// 		if isLikelyImageURL(candidate) {
// 			return candidate
// 		}
// 		return ""
// 	case map[string]any:
// 		if url := strings.TrimSpace(getStringFromMap(v, "url")); url != "" {
// 			if isLikelyImageURL(url) {
// 				return url
// 			}
// 		}
// 		if href := strings.TrimSpace(getStringFromMap(v, "href")); href != "" {
// 			if isLikelyImageURL(href) {
// 				return href
// 			}
// 		}
// 		if icon := v["icon"]; icon != nil {
// 			if nested := extractIconValue(icon); nested != "" {
// 				return nested
// 			}
// 		}
// 		if image := v["image"]; image != nil {
// 			if nested := extractIconValue(image); nested != "" {
// 				return nested
// 			}
// 		}
// 	case []any:
// 		for _, item := range v {
// 			if candidate := extractIconValue(item); candidate != "" {
// 				return candidate
// 			}
// 		}
// 	}
// 	return ""
// }

// func extractActorCandidates(value any) []string {
// 	if value == nil {
// 		return nil
// 	}
// 	candidates := make([]string, 0)
// 	switch v := value.(type) {
// 	case string:
// 		if trimmed := strings.TrimSpace(v); trimmed != "" {
// 			candidates = append(candidates, trimmed)
// 		}
// 	case map[string]any:
// 		if id := strings.TrimSpace(getStringFromMap(v, "id")); id != "" {
// 			candidates = append(candidates, id)
// 		}
// 		if url := strings.TrimSpace(getStringFromMap(v, "url")); url != "" {
// 			candidates = append(candidates, url)
// 		}
// 	case []any:
// 		for _, item := range v {
// 			candidates = append(candidates, extractActorCandidates(item)...)
// 		}
// 	}
// 	return candidates
// }

// func normalizeActivityContent(content string, objectMap map[string]any) string {
// 	trimmed := strings.TrimSpace(content)

// 	if trimmed != "" && looksLikeHTML(trimmed) {
// 		return trimmed
// 	}

// 	if converted := convertSourceToHTML(objectMap["source"]); converted != "" {
// 		return converted
// 	}

// 	if trimmed == "" {
// 		return ""
// 	}

// 	return string(mdUtil.MdToHTML([]byte(trimmed)))
// }

// func convertSourceToHTML(source any) string {
// 	if source == nil {
// 		return ""
// 	}

// 	switch value := source.(type) {
// 	case string:
// 		if strings.TrimSpace(value) == "" {
// 			return ""
// 		}
// 		return string(mdUtil.MdToHTML([]byte(value)))
// 	case map[string]any:
// 		mediaType := strings.ToLower(strings.TrimSpace(getStringFromMap(value, "mediaType")))
// 		if mediaType == "" {
// 			mediaType = strings.ToLower(strings.TrimSpace(getStringFromMap(value, "type")))
// 		}

// 		data := strings.TrimSpace(getStringFromMap(value, "content"))
// 		if data == "" {
// 			data = strings.TrimSpace(getStringFromMap(value, "value"))
// 		}
// 		if data == "" && value["text"] != nil {
// 			data = strings.TrimSpace(extractString(value["text"]))
// 		}
// 		if data == "" {
// 			return ""
// 		}

// 		if strings.Contains(mediaType, "markdown") || !looksLikeHTML(data) {
// 			return string(mdUtil.MdToHTML([]byte(data)))
// 		}
// 		return data
// 	default:
// 		return ""
// 	}
// }

// func looksLikeHTML(value string) bool {
// 	if value == "" {
// 		return false
// 	}
// 	return strings.Contains(value, "<") && strings.Contains(value, ">")
// }

// func isLikelyImageURL(value string) bool {
// 	if value == "" {
// 		return false
// 	}
// 	lower := strings.ToLower(value)
// 	if strings.HasPrefix(lower, "data:image/") {
// 		return true
// 	}
// 	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
// 		extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".bmp", ".avif"}
// 		for _, ext := range extensions {
// 			if strings.Contains(lower, ext) {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// // extractString 将任意值转换为字符串
// func extractString(value any) string {
// 	switch v := value.(type) {
// 	case string:
// 		return strings.TrimSpace(v)
// 	case fmt.Stringer:
// 		return strings.TrimSpace(v.String())
// 	case nil:
// 		return ""
// 	default:
// 		bytes, err := json.Marshal(v)
// 		if err != nil {
// 			return ""
// 		}
// 		return string(bytes)
// 	}
// }

// // mustMarshalStrings 将字符串切片序列化成 JSON 字符串
// func mustMarshalStrings(values []string) string {
// 	if len(values) == 0 {
// 		return "[]"
// 	}
// 	bytes, err := json.Marshal(values)
// 	if err != nil {
// 		return "[]"
// 	}
// 	return string(bytes)
// }

// // derivePreferredUsername 根据 Actor URL 推断用户名
// func derivePreferredUsername(actor string) string {
// 	actor = strings.TrimSpace(actor)
// 	if actor == "" {
// 		return ""
// 	}
// 	actor = strings.TrimSuffix(actor, "/")
// 	parts := strings.Split(actor, "/")
// 	if len(parts) == 0 {
// 		return actor
// 	}
// 	return parts[len(parts)-1]
// }

// // parseRFC3339 尝试解析 RFC3339 或 RFC3339Nano 时间
// func parseRFC3339(value string) (time.Time, error) {
// 	layouts := []string{time.RFC3339Nano, time.RFC3339}
// 	for _, layout := range layouts {
// 		if ts, err := time.Parse(layout, value); err == nil {
// 			return ts, nil
// 		}
// 	}
// 	return time.Time{}, fmt.Errorf("invalid time: %s", value)
// }
