package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
	httpUtil "github.com/lin-snow/ech0/internal/util/http"
)

//==========================================================
//	处理前端的 Actor 搜索请求
//==========================================================

// SearchActorByActorID 根据 Actor ID (URL) 搜索远端 Actor 信息
func (fediverseService *FediverseService) SearchActorByActorID(actorID string) (map[string]any, error) {
	actorID = strings.TrimSpace(actorID)
	if actorID == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	resolvedActorURL, err := resolveActorURL(actorID)
	if err != nil {
		return nil, errors.New(commonModel.GET_ACTOR_ERROR)
	}

	body, err := httpUtil.SendRequest(resolvedActorURL, http.MethodGet, httpUtil.Header{
		Header:  "Accept",
		Content: "application/activity+json",
	}, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("%s", commonModel.GET_ACTOR_ERROR)
	}

	var actor map[string]any
	if err := json.Unmarshal(body, &actor); err != nil {
		return nil, fmt.Errorf("%s", commonModel.GET_ACTOR_ERROR)
	}
	if len(actor) == 0 {
		return nil, errors.New(commonModel.GET_ACTOR_ERROR)
	}

	return actor, nil
}

//==========================================================
//	处理前端的 Follow 请求
//==========================================================

// GetFollowStatus 获取关注状态
func (fediverseService *FediverseService) GetFollowStatus(userID uint, target string) (string, error) {
	// 关注的目标 Actor
	target = strings.TrimSpace(target)
	if target == "" {
		return "", errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}
	// 规范化目标 Actor URL
	target, err := resolveActorURL(target)
	if err != nil {
		return "", err
	}

	// 获取当前用户信息
	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return "", err
	}

	// 检查关注状态
	follow, err := fediverseService.fediverseRepository.GetFollowByUserAndObject(context.Background(), user.ID, target)
	if err != nil {
		return "", err
	}
	// 如果没有关注记录，返回 "none",说明没有关注过
	if follow == nil {
		return model.FollowStatusNone, nil
	}
	// 如果有关注记录，返回当前状态
	return follow.Status, nil
}

// FollowActor 发送关注请求
func (fediverseService *FediverseService) FollowActor(userID uint, req model.FollowActionRequest) (map[string]string, error) {
	// 关注的目标 Actor
	target := strings.TrimSpace(req.TargetActor)
	if target == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}
	// 规范化目标 Actor URL
	target, err := resolveActorURL(target)
	if err != nil {
		return nil, err
	}

	// 获取当前用户信息
	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	// 构建当前用户的 Actor 信息
	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	// 处理当前的 Server URL
	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	published := time.Now().UTC()
	activityID := fmt.Sprintf("%s/activities/%s/follow/%d", serverURL, actor.PreferredUsername, published.UnixNano())

	// 构建 Follow Activity 的 Payload
	payload, err := buildFollowActivityPayload(&actor, target, activityID, published)
	if err != nil {
		return nil, err
	}

	// 获取目标 Actor 的 Inbox URL
	inboxURL, err := fediverseService.fetchRemoteActorInbox(target)
	if err != nil {
		return nil, err
	}

	// 发送 Follow Activity 到目标 Actor 的 Inbox
	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	// 在本地数据库中保存关注关系，状态为 "pending"
	if err := fediverseService.txManager.Run(func(ctx context.Context) error {
		follow := &model.Follow{
			UserID:     user.ID,
			ActorID:    actor.ID,
			ObjectID:   target,
			ActivityID: activityID,
			Status:     model.FollowStatusPending,
		}
		return fediverseService.fediverseRepository.SaveOrUpdateFollow(ctx, follow)
	}); err != nil {
		return nil, err
	}

	// 返回 Activity ID 给前端
	return map[string]string{
		"activityId": activityID,
	}, nil
}

//==========================================================
//	处理前端的 Unfollow 请求
//==========================================================

// UnfollowActor 发送取消关注请求
func (fediverseService *FediverseService) UnfollowActor(userID uint, req model.FollowActionRequest) (map[string]string, error) {
	target := strings.TrimSpace(req.TargetActor)
	if target == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	follow, err := fediverseService.fediverseRepository.GetFollowByUserAndObject(context.Background(), user.ID, target)
	if err != nil {
		return nil, err
	}
	if follow == nil || follow.ActivityID == "" {
		return nil, errors.New(commonModel.FOLLOW_RELATION_MISSING)
	}

	published := time.Now().UTC()
	undoID := fmt.Sprintf("%s/activities/%s/unfollow/%d", serverURL, actor.PreferredUsername, published.UnixNano())

	payload, err := buildUndoFollowActivityPayload(&actor, target, undoID, follow.ActivityID, published)
	if err != nil {
		return nil, err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(target)
	if err != nil {
		return nil, err
	}

	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	if err := fediverseService.txManager.Run(func(ctx context.Context) error {
		return fediverseService.fediverseRepository.DeleteFollow(ctx, follow.ID)
	}); err != nil {
		return nil, err
	}

	return map[string]string{
		"activityId":       undoID,
		"followActivityId": follow.ActivityID,
	}, nil
}

//==========================================================
//	处理前端的 Like 请求
//==========================================================

// LikeObject 发送点赞请求
func (fediverseService *FediverseService) LikeObject(userID uint, req model.LikeActionRequest) (map[string]string, error) {
	targetActor := strings.TrimSpace(req.TargetActor)
	object := strings.TrimSpace(req.Object)
	if targetActor == "" || object == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	likeID := generateDeterministicActivityID(serverURL, actor.PreferredUsername, "like", object)
	published := time.Now().UTC()

	payload, err := buildLikeActivityPayload(&actor, targetActor, object, likeID, published)
	if err != nil {
		return nil, err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(targetActor)
	if err != nil {
		return nil, err
	}

	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	return map[string]string{
		"activityId": likeID,
	}, nil
}

//==========================================================
//	处理前端的 Undo Like 请求
//==========================================================

// UndoLikeObject 发送取消点赞请求
func (fediverseService *FediverseService) UndoLikeObject(userID uint, req model.LikeActionRequest) (map[string]string, error) {
	targetActor := strings.TrimSpace(req.TargetActor)
	object := strings.TrimSpace(req.Object)
	if targetActor == "" || object == "" {
		return nil, errors.New(commonModel.FEDIVERSE_INVALID_INPUT)
	}

	user, err := fediverseService.commonService.CommonGetUserByUserId(userID)
	if err != nil {
		return nil, err
	}

	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return nil, err
	}

	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return nil, err
	}

	likeID := generateDeterministicActivityID(serverURL, actor.PreferredUsername, "like", object)
	published := time.Now().UTC()
	undoID := fmt.Sprintf("%s/activities/%s/undo-like/%d", serverURL, actor.PreferredUsername, published.UnixNano())

	payload, err := buildUndoLikeActivityPayload(&actor, targetActor, object, likeID, undoID, published)
	if err != nil {
		return nil, err
	}

	inboxURL, err := fediverseService.fetchRemoteActorInbox(targetActor)
	if err != nil {
		return nil, err
	}

	if err := httpUtil.PostActivity(payload, inboxURL, actor.ID); err != nil {
		return nil, err
	}

	return map[string]string{
		"activityId":     undoID,
		"likeActivityId": likeID,
	}, nil
}
