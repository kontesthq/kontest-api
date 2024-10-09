package model

import "time"

type Metadata struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

func (k *Metadata) TableName() string {
	return "kontests_metadata"
}

func NewMetadata() *Metadata {
	return &Metadata{
		ID:            "metadata",
		LastUpdatedAt: time.Now(),
	}
}
