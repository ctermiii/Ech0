package model

// FollowActionRequest 前端发起关注/取关请求体
type FollowActionRequest struct {
	TargetActor string `json:"targetActor" binding:"required"`
}

// LikeActionRequest 前端发起点赞/取消点赞请求体
type LikeActionRequest struct {
	TargetActor string `json:"targetActor" binding:"required"`
	Object      string `json:"object" binding:"required"`
	ObjectType  string `json:"objectType,omitempty"`
}
