package utils

import (
	"kontest-api/repository"
	"kontest-api/repository/impl"
	"kontest-api/service"
	"log/slog"
)

// Dependencies holds the application's repositories and services
type Dependencies struct {
	KontestRepository  repository.KontestRepository
	MetadataRepository repository.MetadataRepository
	KontestService     *service.KontestService
}

// NewDependencies initializes the Dependencies struct
func NewDependencies(kontestRepository repository.KontestRepository, metadataRepository repository.MetadataRepository) *Dependencies {
	// print the kontestRepository and metadataRepository
	slog.Info("KontestRepository", slog.Any("kontestRepository", kontestRepository))
	slog.Info("MetadataRepository", slog.Any("metadataRepository", metadataRepository))

	return &Dependencies{
		KontestRepository:  kontestRepository,
		MetadataRepository: metadataRepository,
		KontestService:     service.NewKontestService(kontestRepository, metadataRepository),
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
