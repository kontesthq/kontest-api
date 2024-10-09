package repository

import (
	"kontest-api/model"
	"time"
)

// MetadataRepository defines methods for metadata operations.
type MetadataRepository interface {
	Save(metadata *model.Metadata)
	GetLastUpdatedAt() time.Time
}
