package impl

import (
	"kontest-api/database"
	"kontest-api/model"
	"log"
	"time"
)

// MetadataRepositoryImpl is a concrete implementation of the MetadataRepository interface.
type MetadataRepositoryImpl struct{}

// NewMetadataRepository creates a new instance of MetadataRepositoryImpl.
func NewMetadataRepository() *MetadataRepositoryImpl {
	return &MetadataRepositoryImpl{}
}

// Save saves the metadata to the database.
func (repo *MetadataRepositoryImpl) Save(metadata *model.Metadata) {
	if err := database.GetDB().Save(metadata).Error; err != nil {
		log.Fatalf("Error saving metadata: %v", err)
	}
}

// GetLastUpdatedAt fetches the last updated timestamp from the database.
func (repo *MetadataRepositoryImpl) GetLastUpdatedAt() time.Time {
	var metadata model.Metadata
	if err := database.GetDB().Order("last_updated_at desc").First(&metadata).Error; err != nil {
		log.Printf("Error fetching last updated time: %v", err)
		return time.Time{} // Return zero time if error occurs
	}
	return metadata.LastUpdatedAt
}
