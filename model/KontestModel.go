package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"kontest-api/utils/enums"
)

// KontestModel represents a record in the kontests table
type KontestModel struct {
	ID               uuid.UUID              `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()" json:"id"` // Use UUID type as primary key
	Name             string                 `gorm:"not null" json:"name"`
	URL              string                 `json:"url"`               // Actual URL of the contest
	StartTime        string                 `json:"start_time"`        // Contest start time (stored as string)
	EndTime          string                 `json:"end_time"`          // Contest end time (stored as string)
	Location         string                 `json:"location"`          // Site on which the contest is hosted
	Status           string                 `json:"status"`            // Status of the contest
	SiteAbbreviation enums.SiteAbbreviation `json:"site_abbreviation"` // Abbreviation of the contest site
}

// TableName sets the table name for the KontestModel struct
func (k *KontestModel) TableName() string {
	return "kontests"
}

// BeforeCreate is a GORM hook that runs before inserting a new record into the DB
func (k *KontestModel) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a UUID if the ID is not already set
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}

	// Generate the site abbreviation based on the location using the utility function
	k.SiteAbbreviation = enums.GetAbbreviation(k.Location)
	return nil
}
