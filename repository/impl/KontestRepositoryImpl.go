package impl

import (
	"kontest-api/database"
	"kontest-api/model"
	"log"
)

// KontestRepositoryImpl is a concrete implementation of the KontestRepository interface.
type KontestRepositoryImpl struct{}

// NewKontestRepository creates a new instance of KontestRepositoryImpl.
func NewKontestRepository() *KontestRepositoryImpl {
	return &KontestRepositoryImpl{}
}

// Save saves a contest to the database.
func (repo *KontestRepositoryImpl) Save(kontest model.KontestModel) {
	if err := database.GetDB().Create(&kontest).Error; err != nil {
		// Handle the error appropriately
		log.Fatalf("Error saving contest: %v", err)
	}
}

// DeleteAll deletes all contests from the database.
func (repo *KontestRepositoryImpl) DeleteAll() {
	if err := database.GetDB().Unscoped().Delete(&model.KontestModel{}, "1=1").Error; err != nil {
		// Handle the error appropriately
		log.Fatalf("Error deleting contests: %v", err)
	}
}
