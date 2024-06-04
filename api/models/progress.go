package models

import (
	"time"
)

// 播放进度
type Progress struct {
	Id        uint      `json:"id"`
	UserId    string    `json:"user_id"`
	Data      string    `json:"data"`
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
