package models

import (
	"time"
)

// 剧集进度，第几集
type Progress struct {
	Id        uint      `json:"id"`
	UserId    string    `json:"user_id"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 播放进度，第几秒
type ProgressTv struct {
	Id        uint      `json:"id"`
	UserId    string    `json:"user_id"`
	TvId      uint      `json:"tv_id"`
	SeasonId  uint      `json:"season_id"`
	Time      int       `json:"time"`
	TvPath    string    `json:"tv_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProgressRequestBody struct {
	Data YourModel `json:"data"`
}

type ArtPlayerSettings struct {
	Times map[string]float64 `json:"times"`
}
type YourModel struct {
	ArtPlayerSettings ArtPlayerSettings      `json:"artplayer_settings"`
	TV                map[string]interface{} `json:"tv"`
}

type Request struct {
	Data string `json:"data"`
}
