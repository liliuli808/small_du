package model

import "time"

// BilibiliVideo B站视频表
type BilibiliVideo struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BVID            string    `gorm:"type:varchar(64);uniqueIndex" json:"bvid"`
	AID             int64     `json:"aid"`
	CID             int64     `json:"cid"`
	Title           string    `gorm:"type:text" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	OwnerName       string    `gorm:"type:varchar(255)" json:"owner_name"`
	DurationSeconds int       `json:"duration_seconds"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BilibiliVideo) TableName() string {
	return "bilibili_videos"
}

// VideoPage 视频分P信息
type VideoPage struct {
	CID      int64  `json:"cid"`
	Page     int    `json:"page"`
	Part     string `json:"part"`
	Duration int    `json:"duration"`
}

// VideoInfo 视频完整信息
type VideoInfo struct {
	AID         int64          `json:"aid"`
	BVID        string         `json:"bvid"`
	CID         int64          `json:"cid"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	OwnerName   string         `json:"owner_name"`
	Duration    int            `json:"duration"`
	Pages       []VideoPage    `json:"pages"`
	Subtitles   []SubtitleMeta `json:"subtitles"`
}

// SubtitleMeta 字幕元信息
type SubtitleMeta struct {
	ID          int64  `json:"id"`
	Lan         string `json:"lan"`
	LanDoc      string `json:"lan_doc"`
	SubtitleURL string `json:"subtitle_url"`
	AIStatus    int    `json:"ai_status"`
}

// SubtitleSegment 字幕片段
type SubtitleSegment struct {
	From    float64 `json:"from"`
	To      float64 `json:"to"`
	Content string  `json:"content"`
}

// SubtitleFile 字幕文件结构
type SubtitleFile struct {
	Body []SubtitleSegment `json:"body"`
}

// CommentItem 评论项
type CommentItem struct {
	RPID    int64  `json:"rpid"`
	Message string `json:"message"`
	Like    int    `json:"like"`
	IsTop   bool   `json:"is_top"`
	Member  string `json:"member"`
}
