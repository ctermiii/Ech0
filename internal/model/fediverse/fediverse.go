package model

import (
	"encoding/json"
	"time"
)

// ActivityPubError 定义 ActivityPub 协议的错误响应格式
type ActivityPubError struct {
	Context string `json:"@context"` // ActivityStreams 上下文
	Type    string `json:"type"`     // 类型，固定为 "Error"
	Error   string `json:"error"`    // 错误信息
	Status  int    `json:"status"`   // HTTP 状态码
}

// ActivityType 定义常见的 ActivityPub 活动类型
const (
	ActivityTypeCreate   string = "Create"
	ActivityTypeFollow   string = "Follow"
	ActivityTypeLike     string = "Like"
	ActivityTypeAccept   string = "Accept"
	ActivityTypeAnnounce string = "Announce"
	ActivityTypeUndo     string = "Undo"
)

const (
	DefaultCollectionPageSize = 20
	MaxCollectionPageSize     = 80
)

const (
	FollowStatusNone     = "none"
	FollowStatusPending  = "pending"
	FollowStatusAccepted = "accepted"
	FollowStatusRejected = "rejected"
)

// Link 是 WebFingerResponse 中的链接对象
type Link struct {
	Rel  string `json:"rel"`            // 链接关系，例如 "self"
	Type string `json:"type,omitempty"` // MIME 类型，例如 application/activity+json
	Href string `json:"href"`           // 链接 URL
}

// Activity 表示 ActivityPub 的 Activity
type Activity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"-"`                 // 数据库主键
	Context      any       `gorm:"type:text;not null"       json:"@context"`          // ActivityStreams 上下文，可以是字符串或数组
	ActivityID   string    `gorm:"size:512;unique;not null" json:"id"`                // Activity URL
	Type         string    `gorm:"size:64;not null"         json:"type"`              // Create, Follow, Like, Accept...
	ActorID      string    `gorm:"index;not null"           json:"-"`                 // 关联的用户 ID
	ActorURL     string    `gorm:"size:512;not null"        json:"actor"`             // Actor URL
	Object       any       `gorm:"-"                        json:"object,omitempty"`  // 原始 Object 字段，可能是字符串或对象
	ObjectID     string    `gorm:"size:512;not null"        json:"-"`                 // 目标对象 URL
	ObjectType   string    `gorm:"size:64;not null"         json:"-"`                 // 目标对象类型
	Published    time.Time `                                json:"published"`         // 发布时间
	To           []string  `gorm:"type:text"                json:"to,omitempty"`      // 接收者列表，序列化存储
	Cc           []string  `gorm:"type:text"                json:"cc,omitempty"`      // 补充接收列表
	Summary      string    `gorm:"type:text"                json:"summary,omitempty"` // 可选描述
	ActivityJSON string    `gorm:"type:text;not null"       json:"-"`                 // 原始 Activity JSON
	Delivered    bool      `gorm:"default:false"            json:"-"`                 // 是否投递
	CreatedAt    time.Time `gorm:"autoCreateTime"           json:"-"`                 // 创建时间
}

// Object 内容对象表 (存储 Note, Article, Image 等等)
type Object struct {
	Context      any            `gorm:"-"                        json:"@context,omitempty"`
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"-"`
	ObjectID     string         `gorm:"size:512;unique;not null" json:"id"`                     // 全局唯一 URL
	Type         string         `gorm:"size:64;not null"         json:"type"`                   // Note, Article, Image...
	AttributedTo string         `gorm:"size:512"                 json:"attributedTo,omitempty"` // actor URL
	Content      string         `gorm:"type:text"                json:"content,omitempty"`      // 主要内容
	Source       map[string]any `gorm:"-"                        json:"source,omitempty"`       // 原始内容，可能包含 mediaType 和 content 字段
	Attachments  []Attachment   `gorm:"-"                        json:"attachment,omitempty"`   // 附件 URL 列表，序列化存储
	Published    time.Time      `                                json:"published,omitempty"`
	To           []string       `gorm:"-"                        json:"to,omitempty"` // 序列化成 JSON 存储
	Cc           []string       `gorm:"-"                        json:"cc,omitempty"` // 同上
	ObjectJSON   string         `gorm:"type:text"                json:"-"`            // 完整 JSON，便于恢复
	CreatedAt    time.Time      `gorm:"autoCreateTime"           json:"-"`
}

// Attachment 附件对象
type Attachment struct {
	Type      string   `json:"type"`               // "Image"、"Video" 等
	MediaType string   `json:"mediaType"`          // MIME 类型
	URL       string   `json:"url"`                // 媒体 URL
	Name      string   `json:"name,omitempty"`     // 媒体名称
	Caption   string   `json:"caption,omitempty"`  // 媒体说明
	Width     int      `json:"width,omitempty"`    // 宽度
	Height    int      `json:"height,omitempty"`   // 高度
	Duration  string   `json:"duration,omitempty"` // 视频/音频时长
	Preview   *Preview `json:"preview,omitempty"`  // 预览信息
}

// Preview 预览对象
type Preview struct {
	Type      string `json:"type"`
	MediaType string `json:"mediaType"`
	URL       string `json:"url"`
}

// PublicKey 公钥信息
type PublicKey struct {
	ID           string `json:"id"`           // 公钥 ID，通常为 Actor ID + "#main-key"
	Owner        string `json:"owner"`        // 所有者，通常为 Actor ID
	PublicKeyPem string `json:"publicKeyPem"` // PEM 格式的公钥
	Type         string `json:"type"`         // "PublicKey"
}

// Actor ActivityPub Actor 信息
type Actor struct {
	Context           []interface{} `json:"@context"`          // 上下文，可以是字符串或对象的数组
	ID                string        `json:"id"`                // Actor 的唯一标识 URL，格式通常为 http(s)://domain/users/username
	Type              string        `json:"type"`              // Actor 类型，通常为 "Person"
	Name              string        `json:"name"`              // 显示名称
	PreferredUsername string        `json:"preferredUsername"` // 用户名
	Summary           string        `json:"summary"`           // 简短介绍
	Icon              Preview       `json:"icon,omitempty"`    // 头像信息
	Image             Preview       `json:"image,omitempty"`   // 封面图片
	Followers         string        `json:"followers"`         // 粉丝列表 URL
	Following         string        `json:"following"`         // 关注列表 URL
	Inbox             string        `json:"inbox"`             // 收件箱 URL
	Outbox            string        `json:"outbox"`            // 发件箱 URL
	PublicKey         PublicKey     `json:"publicKey"`         // 公钥信息
}

// Follow 表：存储关注请求及状态
type Follow struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint      `gorm:"not null;index"           json:"user_id"`     // 发起关注的用户数据库 ID
	ActorID    string    `gorm:"size:512;not null;index"  json:"actor_id"`    // 发起关注的 Actor URL, 格式通常为 http(s)://domain/users/username
	ObjectID   string    `gorm:"size:512;not null;index"  json:"object_id"`   // 被关注的 Actor URL, 格式通常为 http(s)://domain/users/username
	ActivityID string    `gorm:"size:512"                 json:"activity_id"` // Follow 活动 ID，便于撤销
	Status     string    `gorm:"size:20;not null"         json:"status"`      // pending, accepted, rejected
	CreatedAt  time.Time `gorm:"autoCreateTime"           json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"           json:"updated_at"`
}

// Follower 表：存储已接受的关注关系
type Follower struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActorID   string    `gorm:"size:512;not null;index"  json:"actor_id"` // 粉丝 Actor URL, 格式通常为 http(s)://domain/users/username
	UserID    uint      `gorm:"not null;index"           json:"user_id"`  // 被关注用户的数据库 ID
	CreatedAt time.Time `gorm:"autoCreateTime"           json:"created_at"`
}

// InboxStatus 收件箱中存储的远端推文记录，供后续时间线展示使用
type InboxStatus struct {
	ID                     uint      `gorm:"primaryKey;autoIncrement"                             json:"id"`
	UserID                 uint      `gorm:"not null;index:idx_user_activity,priority:1"          json:"user_id"` // 本地用户 ID
	ActivityID             string    `gorm:"size:512;not null;index:idx_user_activity,priority:2" json:"activity_id"`
	ActorID                string    `gorm:"size:512;not null"                                    json:"actor_id"`                 // 发送者 Actor URL
	ActorPreferredUsername string    `gorm:"size:128"                                             json:"actor_preferred_username"` // 发送者用户名（推断）
	ActorDisplayName       string    `gorm:"size:255"                                             json:"actor_display_name"`       // 发送者显示名称
	ActorAvatar            string    `gorm:"size:512"                                             json:"actor_avatar"`             // 发送者头像 URL
	ObjectID               string    `gorm:"size:512;not null"                                    json:"object_id"`                // Object 的唯一 URL
	ObjectType             string    `gorm:"size:64"                                              json:"object_type"`              // Object 类型，例如 Note
	ObjectAttributedTo     string    `gorm:"size:512"                                             json:"object_attributed_to"`     // Object 的 attributedTo 字段
	Summary                string    `gorm:"type:text"                                            json:"summary"`                  // Activity 或 Object 中的摘要
	Content                string    `gorm:"type:text"                                            json:"content"`                  // Object 内容，通常为 HTML
	To                     string    `gorm:"type:text"                                            json:"to"`                       // 接收者列表，JSON 字符串
	Cc                     string    `gorm:"type:text"                                            json:"cc"`                       // 抄送列表，JSON 字符串
	RawActivity            string    `gorm:"type:text;not null"                                   json:"raw_activity"`             // 完整的 Activity JSON
	RawObject              string    `gorm:"type:text"                                            json:"raw_object"`               // 完整的 Object JSON
	PublishedAt            time.Time `gorm:"index"                                                json:"published_at"`             // 消息发布时间
	CreatedAt              time.Time `gorm:"autoCreateTime"                                       json:"created_at"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"                                       json:"updated_at"`
}

// TimelineItem 联邦时间线条目响应体
type TimelineItem struct {
	ID                     uint            `json:"id"`                     // 本地收件箱记录 ID
	ActivityID             string          `json:"activityId"`             // Activity ID
	ActorID                string          `json:"actorId"`                // 发送者 Actor URL
	ActorPreferredUsername string          `json:"actorPreferredUsername"` // 发送者用户名
	ActorDisplayName       string          `json:"actorDisplayName"`       // 发送者显示名称
	ActorAvatar            string          `json:"actorAvatar"`            // 发送者头像
	ObjectID               string          `json:"objectId"`               // Object 全局 URL
	ObjectType             string          `json:"objectType"`             // Object 类型
	ObjectAttributedTo     string          `json:"objectAttributedTo"`     // Object attributedTo
	Summary                string          `json:"summary"`                // 摘要
	Content                string          `json:"content"`                // HTML 内容
	To                     []string        `json:"to"`                     // 收件人列表
	Cc                     []string        `json:"cc"`                     // 抄送列表
	RawActivity            json.RawMessage `json:"rawActivity,omitempty"`  // 完整 Activity JSON
	RawObject              json.RawMessage `json:"rawObject,omitempty"`    // 完整 Object JSON
	PublishedAt            time.Time       `json:"publishedAt"`            // 发布时间
	CreatedAt              time.Time       `json:"createdAt"`              // 本地记录创建时间
	UpdatedAt              time.Time       `json:"updatedAt"`              // 本地记录更新时间
}
