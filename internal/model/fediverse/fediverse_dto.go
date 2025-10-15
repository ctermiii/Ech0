package model

// FollowActionRequest 前端发起关注/取关请求体
type FollowActionRequest struct {
	TargetActor string `json:"targetActor" binding:"required"`
}

// LikeActionRequest 前端发起点赞/取消点赞请求体
type LikeActionRequest struct {
	TargetActor string `json:"targetActor"          binding:"required"`
	Object      string `json:"object"               binding:"required"`
	ObjectType  string `json:"objectType,omitempty"`
}

// WebFingerResponse 是 WebFinger 返回的标准结构
type WebFingerResponse struct {
	Subject string   `json:"subject"`           // 用户标识，例如 acct:alice@domain.com
	Aliases []string `json:"aliases,omitempty"` // 可选：用户的别名 URL
	Links   []Link   `json:"links"`             // 与用户相关的资源链接
}

// OutboxResponse 定义 Outbox 的响应格式
type OutboxResponse struct {
	Context    interface{} `json:"@context"` // 可以是字符串或数组
	ID         string      `json:"id"`
	Type       string      `json:"type"` // "OrderedCollection"
	TotalItems int         `json:"totalItems"`
	First      string      `json:"first,omitempty"`
	Last       string      `json:"last,omitempty"`
}

// OutboxPage 表示分页形式的 Outbox
type OutboxPage struct {
	Context      any        `json:"@context,omitempty"`
	ID           string     `json:"id"`
	Type         string     `json:"type"` // "OrderedCollectionPage"
	PartOf       string     `json:"partOf"`
	Next         string     `json:"next,omitempty"`
	Prev         string     `json:"prev,omitempty"`
	OrderedItems []Activity `json:"orderedItems"`
}

// FollowersResponse 跟 OutboxResponse 类似
type FollowersResponse struct {
	Context    any    `json:"@context"`
	ID         string `json:"id"`
	Type       string `json:"type"` // "OrderedCollection"
	TotalItems int    `json:"totalItems"`
	First      any    `json:"first,omitempty"`
	// 如果不分页，可以直接用
	OrderedItems []string `json:"orderedItems,omitempty"` // 里面是 follower 的 Actor URL
}

// FollowersPage 如果你要分页的话
type FollowersPage struct {
	Context      any      `json:"@context,omitempty"`
	ID           string   `json:"id"`
	Type         string   `json:"type"` // "OrderedCollectionPage"
	PartOf       string   `json:"partOf"`
	Next         string   `json:"next,omitempty"`
	Prev         string   `json:"prev,omitempty"`
	OrderedItems []string `json:"orderedItems"`
}

// FollowingResponse 跟 FollowersResponse 类似
type FollowingResponse struct {
	Context    any    `json:"@context"`
	ID         string `json:"id"`
	Type       string `json:"type"` // "OrderedCollection"
	TotalItems int    `json:"totalItems"`
	First      any    `json:"first,omitempty"`
	// 如果不分页，可以直接用
	OrderedItems []string `json:"orderedItems,omitempty"` // 里面是 following 的 Actor URL
}

// FollowingPage 如果你要分页的话
type FollowingPage struct {
	Context      any      `json:"@context,omitempty"`
	ID           string   `json:"id"`
	Type         string   `json:"type"` // "OrderedCollectionPage"
	PartOf       string   `json:"partOf"`
	Next         string   `json:"next,omitempty"`
	Prev         string   `json:"prev,omitempty"`
	OrderedItems []string `json:"orderedItems"`
}
