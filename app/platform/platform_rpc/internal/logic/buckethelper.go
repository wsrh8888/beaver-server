package logic

import (
	"beaver/app/platform/platform_models"

	"gorm.io/gorm"
)

func bucketNameMap(db *gorm.DB, bucketIDs []string) map[string]string {
	result := make(map[string]string, len(bucketIDs))
	if len(bucketIDs) == 0 {
		return result
	}

	uniqueIDs := make([]string, 0, len(bucketIDs))
	seen := make(map[string]struct{}, len(bucketIDs))
	for _, id := range bucketIDs {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	if len(uniqueIDs) == 0 {
		return result
	}

	var buckets []platform_models.TrackBucket
	db.Where("bucket_id IN ?", uniqueIDs).Find(&buckets)
	for _, bucket := range buckets {
		result[bucket.BucketID] = bucket.Name
	}
	return result
}

func bucketNameByID(db *gorm.DB, bucketID string) string {
	if bucketID == "" {
		return ""
	}
	var bucket platform_models.TrackBucket
	if err := db.Where("bucket_id = ?", bucketID).Take(&bucket).Error; err != nil {
		return ""
	}
	return bucket.Name
}
