package models

import "time"

type ZoomMeet struct {
	ID        string     `json:"id" gorm:"primary_key;not null;type:varchar(36);index"`
	UserID    string     `json:"user_id" gorm:"not null;type:varchar(250);index"`
	MeetingID int64      `json:"meeting_id" gorm:"not null;type:int"`
	Topic     string     `json:"topic" gorm:"type:text"`
	StartTime time.Time  `json:"start_time" gorm:"not null"`
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
	UserID    string    `json:"user_id"`
	Topic     string    `json:"topic"`
	StartTime time.Time `json:"start_time"`
	Duration  int       `json:"duration"`
	Status    Status    `json:"status"`
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
	UserID    string     `json:"user_id"`
	MeetingID int64      `json:"meeting_id" gorm:"not null;type:int"`
	Topic     string     `json:"topic" gorm:"type:text"`
	StartTime time.Time  `json:"start_time" gorm:"not null"`
	Duration  int        `json:"duration" gorm:"type:int"`
	JoinURL   string     `json:"join_url" gorm:"type:text"`
	Status    Status     `json:"status" gorm:"not null;type:varchar(10);index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type ZoomMeetUpdate struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	MeetingID int64     `json:"meeting_id" gorm:"not null;type:int"`
	Topic     string    `json:"topic" gorm:"type:text"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	Duration  int       `json:"duration" gorm:"type:int"`
	JoinURL   string    `json:"join_url" gorm:"type:text"`
	Status    Status    `json:"status" gorm:"not null;type:varchar(10);index"`
}
