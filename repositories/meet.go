package repositories

import (
	"backend-b7/models"
	"backend-b7/pkg/utils"
	"fmt"

	"gorm.io/gorm"
)

type meetRepository struct {
	db *gorm.DB
}

type MeetRepositoryInterface interface {
	CreateMeet(meet *models.ZoomMeet) error
	GetMeets(pagination utils.Pagination, where map[string]string) ([]models.ZoomMeetView, int64, error)
	GetMeetById(id string) (models.ZoomMeetView, error)
	UpdateMeet(meet *models.ZoomMeetUpdate) error
	DeleteMeet(meet *models.ZoomMeetUpdate) (err error)
}

func NewMeetRepository(db *gorm.DB) MeetRepositoryInterface {
	return &meetRepository{
		db: db,
	}
}

func (pr *meetRepository) CreateMeet(meet *models.ZoomMeet) error {
	return pr.db.Create(meet).Error
}

func (pr *meetRepository) GetMeets(pagination utils.Pagination, where map[string]string) ([]models.ZoomMeetView, int64, error) {
	var count int64
	var err error
	var sortField, sortDirection string
	var meets []models.ZoomMeetView

	queryBuilder := pr.db.
		Table("meets").Select("meets.*").
		Where("meets.status <> ?", models.StatusDeleted)

	if id, ok := where["id"]; ok && id != "" {
		queryBuilder = queryBuilder.Where(`meets.id = ?`, id)
	}

	if meetingID, ok := where["meeting_id"]; ok && meetingID != "" {
		queryBuilder = queryBuilder.Where(`meets.meeting_id = ?`, meetingID)
	}

	if topic, ok := where["topic"]; ok && topic != "" {
		topic := fmt.Sprintf("%%%s%%", topic)
		queryBuilder = queryBuilder.Where(`meets."topic" = ?`, topic)
	}

	if pagination.SortField != "" {
		if pagination.SortField == "topic" {
			sortField = `INITCAP(meets."topic")`
		} else if pagination.SortField == "duration" {
			sortField = "meets.duration"
		} else if pagination.SortField == "start_time" {
			sortField = "meets.start_time"
		} else {
			sortField = "meets.created_at"
		}
	} else {
		sortField = "meets.created_at"
	}

	if pagination.SortDirection != "" {
		sortDirection = pagination.SortDirection
	} else {
		sortDirection = models.SortDirectionDESC.String()
	}

	err = queryBuilder.Count(&count).Error
	if err != nil {
		return nil, count, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	orderBy := fmt.Sprintf("%s %s", sortField, sortDirection)
	limitBuilder := queryBuilder.Limit(pagination.Limit).Offset(offset).Order(orderBy)

	result := limitBuilder.Scan(&meets)
	if result.Error != nil {
		return nil, count, result.Error
	}

	return meets, count, nil
}

func (pr *meetRepository) GetMeetById(id string) (models.ZoomMeetView, error) {
	var meet models.ZoomMeetView
	queryBuilder := pr.db.
		Table("meets").Select("meets.*").
		Where("meets.status <> ?", models.StatusDeleted)
	if err := queryBuilder.Where("meets.id = ?", id).Scan(&meet).Error; err != nil {
		return meet, err
	}
	return meet, nil
}

func (pr *meetRepository) UpdateMeet(meet *models.ZoomMeetUpdate) error {
	return pr.db.
		Model(&models.ZoomMeet{ID: meet.ID}).
		Updates(
			map[string]interface{}{
				"meeting_id": meet.MeetingID,
				"topic":      meet.Topic,
				"start_time": meet.StartTime,
				"duration":   meet.Duration,
				"join_url":   meet.Status,
				"status":     meet.Status,
				"updated_at": gorm.Expr("now()"),
			},
		).Error
}

func (pr *meetRepository) DeleteMeet(meet *models.ZoomMeetUpdate) (err error) {
	return pr.db.
		Model(&models.ZoomMeet{ID: meet.ID}).
		Updates(
			map[string]interface{}{
				"status":     models.StatusDeleted.String(),
				"updated_at": gorm.Expr("now()"),
			},
		).Error
}
