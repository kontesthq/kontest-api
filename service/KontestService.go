package service

import (
	"kontest-api/database"
	"kontest-api/model" // Update with your actual model package path
)

// KontestService interface defines the methods for contest operations
type KontestService interface {
	GetContests(page, perPage int) ([]map[string]string, error)
	GetContestsOfSpecificSites(sites []string, page, perPage int) ([]map[string]string, error)
}

// kontestService is the implementation of KontestService
type kontestService struct {
}

// NewKontestService creates a new instance of KontestService
func NewKontestService() KontestService {
	return &kontestService{}
}

// GetContests retrieves a paginated list of contests
func (s *kontestService) GetContests(page, perPage int) ([]map[string]string, error) {
	var contests []model.KontestModel

	// Calculate offset for pagination
	offset := (page - 1) * perPage

	// Query the database for contests
	if err := database.GetDB().Offset(offset).Limit(perPage).Find(&contests).Error; err != nil {
		return nil, err
	}

	// Convert to a slice of maps for easier JSON serialization
	result := make([]map[string]string, len(contests))
	for i, contest := range contests {
		result[i] = map[string]string{
			"id":                contest.ID.String(),
			"name":              contest.Name,
			"url":               contest.URL,
			"start_time":        contest.StartTime,
			"end_time":          contest.EndTime,
			"location":          contest.Location,
			"status":            contest.Status,
			"site_abbreviation": string(contest.SiteAbbreviation), // Assuming it's a string type
		}
	}

	return result, nil
}

// GetContestsOfSpecificSites retrieves contests for specific sites with pagination
func (s *kontestService) GetContestsOfSpecificSites(sites []string, page, perPage int) ([]map[string]string, error) {
	var contests []model.KontestModel

	// Calculate offset for pagination
	offset := (page - 1) * perPage

	// Query the database for contests based on specific sites
	if err := database.GetDB().Where("site_abbreviation IN ?", sites).Offset(offset).Limit(perPage).Find(&contests).Error; err != nil {
		return nil, err
	}

	// Convert to a slice of maps for easier JSON serialization
	result := make([]map[string]string, len(contests))
	for i, contest := range contests {
		result[i] = map[string]string{
			"id":                contest.ID.String(),
			"name":              contest.Name,
			"url":               contest.URL,
			"start_time":        contest.StartTime,
			"end_time":          contest.EndTime,
			"location":          contest.Location,
			"status":            contest.Status,
			"site_abbreviation": string(contest.SiteAbbreviation), // Assuming it's a string type
		}
	}

	return result, nil
}
