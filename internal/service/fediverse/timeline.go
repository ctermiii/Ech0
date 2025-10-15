package service

// import (
// 	"context"
// 	"encoding/json"
// 	"strings"

// 	commonModel "github.com/lin-snow/ech0/internal/model/common"
// 	model "github.com/lin-snow/ech0/internal/model/fediverse"
// )

// // GetTimeline 获取当前用户关注的远端推文时间线
// func (fediverseService *FediverseService) GetTimeline(
// 	userID uint,
// 	page, pageSize int,
// ) (commonModel.PageQueryResult[[]model.TimelineItem], error) {
// 	page, pageSize = normalizePageParams(page, pageSize)

// 	statuses, total, err := fediverseService.fediverseRepository.ListInboxStatuses(
// 		context.Background(),
// 		userID,
// 		page,
// 		pageSize,
// 	)
// 	if err != nil {
// 		return commonModel.PageQueryResult[[]model.TimelineItem]{}, err
// 	}

// 	items := make([]model.TimelineItem, 0, len(statuses))
// 	for _, status := range statuses {
// 		items = append(items, convertInboxStatusToTimeline(status))
// 	}

// 	return commonModel.PageQueryResult[[]model.TimelineItem]{
// 		Total: total,
// 		Items: items,
// 	}, nil
// }

// func convertInboxStatusToTimeline(status model.InboxStatus) model.TimelineItem {
// 	return model.TimelineItem{
// 		ID:                     status.ID,
// 		ActivityID:             status.ActivityID,
// 		ActorID:                status.ActorID,
// 		ActorPreferredUsername: status.ActorPreferredUsername,
// 		ActorDisplayName:       status.ActorDisplayName,
// 		ActorAvatar:            status.ActorAvatar,
// 		ObjectID:               status.ObjectID,
// 		ObjectType:             status.ObjectType,
// 		ObjectAttributedTo:     status.ObjectAttributedTo,
// 		Summary:                status.Summary,
// 		Content:                status.Content,
// 		To:                     parseRecipients(status.To),
// 		Cc:                     parseRecipients(status.Cc),
// 		RawActivity:            rawJSONOrNil(status.RawActivity),
// 		RawObject:              rawJSONOrNil(status.RawObject),
// 		PublishedAt:            status.PublishedAt,
// 		CreatedAt:              status.CreatedAt,
// 		UpdatedAt:              status.UpdatedAt,
// 	}
// }

// func parseRecipients(raw string) []string {
// 	raw = strings.TrimSpace(raw)
// 	if raw == "" {
// 		return []string{}
// 	}

// 	var recipients []string
// 	if err := json.Unmarshal([]byte(raw), &recipients); err != nil {
// 		return []string{}
// 	}

// 	return recipients
// }

// func rawJSONOrNil(raw string) json.RawMessage {
// 	raw = strings.TrimSpace(raw)
// 	if raw == "" {
// 		return nil
// 	}
// 	return json.RawMessage([]byte(raw))
// }
