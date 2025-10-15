package service

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"gorm.io/gorm"

// 	model "github.com/lin-snow/ech0/internal/model/fediverse"
// 	userModel "github.com/lin-snow/ech0/internal/model/user"
// )

// // handleAcceptActivity 处理远端返回的 Accept 活动，将本地关注状态标记为已接受
// func (fediverseService *FediverseService) handleAcceptActivity(user *userModel.User, activity *model.Activity) error {
// 	if user == nil {
// 		return errors.New("user is nil")
// 	}
// 	if activity == nil {
// 		return errors.New("activity is nil")
// 	}

// 	followActivityID := extractFollowActivityIDFromAccept(activity.Object)
// 	if followActivityID == "" {
// 		followActivityID = strings.TrimSpace(activity.ObjectID)
// 	}
// 	if followActivityID == "" {
// 		return errors.New("accept activity missing follow id")
// 	}

// 	return fediverseService.txManager.Run(func(ctx context.Context) error {
// 		err := fediverseService.fediverseRepository.UpdateFollowStatusByActivityID(
// 			ctx,
// 			user.ID,
// 			followActivityID,
// 			model.FollowStatusAccepted,
// 		)
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				fmt.Printf(
// 					"accept activity references unknown follow: user=%d activity=%s\n",
// 					user.ID,
// 					followActivityID,
// 				)
// 				return nil
// 			}
// 		}
// 		return err
// 	})
// }

// // extractFollowActivityIDFromAccept 从 Accept 活动的 object 中提取原始 Follow Activity ID
// func extractFollowActivityIDFromAccept(object any) string {
// 	switch value := object.(type) {
// 	case string:
// 		return strings.TrimSpace(value)
// 	case map[string]any:
// 		if id := strings.TrimSpace(getStringFromMap(value, "id")); id != "" {
// 			return id
// 		}
// 		if nested, ok := value["object"]; ok {
// 			switch n := nested.(type) {
// 			case string:
// 				s := strings.TrimSpace(n)
// 				if strings.Contains(s, "/activities/") {
// 					return s
// 				}
// 			case map[string]any:
// 				if id := strings.TrimSpace(getStringFromMap(n, "id")); id != "" {
// 					return id
// 				}
// 			}
// 		}
// 	case []any:
// 		for _, item := range value {
// 			if id := extractFollowActivityIDFromAccept(item); id != "" {
// 				return id
// 			}
// 		}
// 	}
// 	return ""
// }
