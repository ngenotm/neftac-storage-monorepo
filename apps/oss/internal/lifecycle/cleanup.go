package lifecycle

import (
	"log"
	"time"
	"neftac/storage/internal/db"
	"neftac/storage/internal/storage"
)

func StartCleanup(days int) {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		cutoff := time.Now().AddDate(0, 0, -days).Unix()
		var old []db.Object
		db.DB.Where("created_at < ? AND is_latest = ?", cutoff, false).Find(&old)
		for _, o := range old {
			var b db.Bucket
			db.DB.Where("id = ?", o.BucketID).First(&b)
			storage.DeleteObject(b.Name, o.Key)
			log.Printf("Expired: %s/%s", b.Name, o.Key)
		}
	}
}
