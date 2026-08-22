package circle

import (
	"beaver/app/circle/circle_models"

	"gorm.io/gorm"
)

type circleCountRow struct {
	CircleID string
	Count    int64
}

func countMembers(db *gorm.DB, circleID string) int64 {
	var count int64
	db.Model(&circle_models.CircleMemberModel{}).Where("circle_id = ?", circleID).Count(&count)
	return count
}

func countPosts(db *gorm.DB, circleID string) int64 {
	var count int64
	db.Model(&circle_models.CirclePostModel{}).
		Where("circle_id = ? AND is_deleted = false", circleID).
		Count(&count)
	return count
}

func countMembersByCircleIDs(db *gorm.DB, circleIDs []string) map[string]int64 {
	result := make(map[string]int64, len(circleIDs))
	if len(circleIDs) == 0 {
		return result
	}
	var rows []circleCountRow
	db.Model(&circle_models.CircleMemberModel{}).
		Select("circle_id, count(*) as count").
		Where("circle_id IN ?", circleIDs).
		Group("circle_id").
		Scan(&rows)
	for _, row := range rows {
		result[row.CircleID] = row.Count
	}
	return result
}

func countPostsByCircleIDs(db *gorm.DB, circleIDs []string) map[string]int64 {
	result := make(map[string]int64, len(circleIDs))
	if len(circleIDs) == 0 {
		return result
	}
	var rows []circleCountRow
	db.Model(&circle_models.CirclePostModel{}).
		Select("circle_id, count(*) as count").
		Where("circle_id IN ? AND is_deleted = false", circleIDs).
		Group("circle_id").
		Scan(&rows)
	for _, row := range rows {
		result[row.CircleID] = row.Count
	}
	return result
}
