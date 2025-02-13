package models

import "time"

type ZoomMeet struct {
	ID        string     `json:"id" gorm:"primary_key;not null;type:varchar(36);index"`
	MeetingID string     `json:"meeting_id" gorm:"not null;type:varchar(250)"`
	Topic     string     `json:"topic" gorm:"type:text"`
	StartTime string     `json:"start_time" gorm:"type:varchar(250)"`
	Duration  int        `json:"duration" gorm:"type:int"`
	JoinURL   string     `json:"join_url" gorm:"type:text"`
	Status    Status     `json:"status" gorm:"not null;type:varchar(10);index"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"default:null"`
}

func (ZoomMeet) TableName() string {
	return "meets"
}

type ZoomMeetRegister struct {
	Topic     string `json:"topic"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
	Status    Status `json:"status"`
}

type ListZoomMeet struct {
	Page      int            `json:"page"`
	Limit     int            `json:"limit"`
	Total     int            `json:"total"`
	TotalPage int            `json:"totalPage"`
	ZoomMeets []ZoomMeetView `json:"zoom_meets"`
}

type ZoomMeetView struct {
	ID        string     `json:"id"`
	MeetingID string     `json:"meeting_id" gorm:"not null;type:varchar(250)"`
	Topic     string     `json:"topic" gorm:"type:text"`
	StartTime string     `json:"start_time" gorm:"type:varchar(250)"`
	Duration  int        `json:"duration" gorm:"type:int"`
	JoinURL   string     `json:"join_url" gorm:"type:text"`
	Status    Status     `json:"status" gorm:"not null;type:varchar(10);index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ZoomMeetUpdate struct {
	ID        string `json:"id"`
	MeetingID string `json:"meeting_id" gorm:"not null;type:varchar(250)"`
	Topic     string `json:"topic" gorm:"type:text"`
	StartTime string `json:"start_time" gorm:"type:varchar(250)"`
	Duration  int    `json:"duration" gorm:"type:int"`
	JoinURL   string `json:"join_url" gorm:"type:text"`
	Status    Status `json:"status" gorm:"not null;type:varchar(10);index"`
}
