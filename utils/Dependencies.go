package utils

import (
	"kontest-api/repository"
	"kontest-api/repository/impl"
)

// Dependencies holds the application's repositories
type Dependencies struct {
	KontestRepository  repository.KontestRepository
	MetadataRepository repository.MetadataRepository
}

// NewDependencies initializes the Dependencies struct
func NewDependencies(kontestRepository repository.KontestRepository, metadataRepository repository.MetadataRepository) *Dependencies {
	return &Dependencies{
		KontestRepository:  kontestRepository,
		MetadataRepository: metadataRepository,
	}
}

// Global variable to hold the application dependencies
var dependencies *Dependencies

// InitializeDependencies sets the global dependencies
func InitializeDependencies() {
	dependencies = NewDependencies(impl.NewKontestRepository(), impl.NewMetadataRepository())
}

// GetDependencies returns the global dependencies
func GetDependencies() *Dependencies {
	return dependencies
}
