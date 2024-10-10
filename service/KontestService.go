// service.go
package service

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"kontest-api/database"
	"kontest-api/model"
	"kontest-api/repository"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const updateInterval = time.Hour

type KontestService struct {
	kontestRepo   repository.KontestRepository
	metadataRepo  repository.MetadataRepository
	url           string
	lastUpdatedAt time.Time
	updateMutex   sync.Mutex
	kontestsCache []model.KontestModel // Cache variable
	isUpdating    sync.Mutex
}

func NewKontestService(kontestRepository repository.KontestRepository, metadataRepository repository.MetadataRepository) *KontestService {
	var kontests []model.KontestModel

	// Fetch contests from the database
	database.GetDB().Find(&kontests)

	return &KontestService{
		kontestRepo:   kontestRepository,
		metadataRepo:  metadataRepository,
		url:           "https://clist.by",
		lastUpdatedAt: metadataRepository.GetLastUpdatedAt(),
		kontestsCache: kontests, // Initialize the cache with fetched contests
	}
}

func (s *KontestService) fetchHtml() {
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Check if an update is needed
	if time.Since(s.lastUpdatedAt) < updateInterval {
		log.Println("Update is not required.")
		return
	}

	// Fetch HTML content from the URL
	resp, err := http.Get(s.url)
	if err != nil {
		log.Fatalf("Failed to fetch HTML content: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch HTML content: received status %s", resp.Status)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body) // Read all the response body
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
		return
	}

	log.Println("Fetched HTML content successfully.")

	kontests := s.parseContests(string(body))

	const timeLayout = "January 2, 2006 15:04:05"

	// Sort the kontests slice
	sort.Slice(kontests, func(i, j int) bool {
		// Parse StartTime
		startTimeI, errI := time.Parse(timeLayout, kontests[i].StartTime)
		startTimeJ, errJ := time.Parse(timeLayout, kontests[j].StartTime)

		// If parsing fails, you might want to handle the error or set a default time
		if errI != nil || errJ != nil {
			// Handle error, e.g., log it or ignore the entry
			return false // or any logic you want to apply when there's an error
		}

		// Sort by StartTime
		if !startTimeI.Equal(startTimeJ) {
			return startTimeI.Before(startTimeJ)
		}

		// Parse EndTime
		endTimeI, errI := time.Parse(timeLayout, kontests[i].EndTime)
		endTimeJ, errJ := time.Parse(timeLayout, kontests[j].EndTime)

		// Handle error similarly
		if errI != nil || errJ != nil {
			return false // or any logic you want to apply when there's an error
		}

		// Sort by EndTime
		if !endTimeI.Equal(endTimeJ) {
			return endTimeI.Before(endTimeJ)
		}

		// Finally sort by SiteAbbreviation
		return kontests[i].SiteAbbreviation < kontests[j].SiteAbbreviation
	})

	// Clear existing contests and save new ones
	s.kontestRepo.DeleteAll()
	s.saveKontestsToDB(kontests)

	// Update cache
	s.kontestsCache = kontests

	// Update metadata
	s.metadataRepo.Save(model.NewMetadata())
	s.lastUpdatedAt = time.Now()
}

func (s *KontestService) parseContests(html string) []model.KontestModel {
	var kontestModels []model.KontestModel

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	k := doc.Find("tr.contest")
	fmt.Println("Total contest rows found:", k.Length())

	k.Each(func(i int, s *goquery.Selection) {
		// Extracting specific fields for each contest
		startTime := s.Find("td.start-time a").Text()
		duration := s.Find("td.duration").Text()
		name := s.Find("td.event a.title-search").Text()
		desc, exists := s.Find("a.data-ace").Attr("data-ace")

		println(startTime, duration, name, desc, exists)

		if !exists {
			log.Println("data-ace attribute not found for contest:", name)
			return
		}

		// Unmarshal the JSON data (handling HTML entities)
		desc = strings.ReplaceAll(desc, "&quot;", "\"")
		var dataAce struct {
			Title    string `json:"title"`
			Desc     string `json:"desc"`
			Location string `json:"location"`
			Time     struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"time"`
		}

		if err := json.Unmarshal([]byte(desc), &dataAce); err != nil {
			log.Printf("Failed to unmarshal JSON: %v\n", err)
			return
		}

		// Extract URL from the desc
		url := strings.Replace(dataAce.Desc, "url: ", "", 1)

		kontest := model.NewKontestModel(name, url, dataAce.Time.Start, dataAce.Time.End, dataAce.Location, "")

		kontestModels = append(kontestModels, *kontest)
	})

	log.Println("Parsed contests successfully.")
	return kontestModels
}

func (s *KontestService) saveKontestsToDB(kontests []model.KontestModel) {
	for _, kontest := range kontests {
		s.kontestRepo.Save(kontest)
	}
}

// GetContests retrieves a paginated list of contests
func (s *KontestService) GetContests(page, perPage int) ([]map[string]string, error) {
	s.fetchHtmlIfNeeded()

	// Calculate offset for pagination
	offset := (page - 1) * perPage

	// Limit the result size to perPage
	if offset >= len(s.kontestsCache) {
		return []map[string]string{}, nil // Return an empty result if the offset is out of range
	}

	// Convert to a slice of maps for easier JSON serialization
	result := make([]map[string]string, 0, perPage) // Initialize with capacity of perPage
	for i := offset; i < offset+perPage && i < len(s.kontestsCache); i++ {
		contest := s.kontestsCache[i]
		result = append(result, map[string]string{
			"id":                contest.ID.String(),
			"name":              contest.Name,
			"url":               contest.URL,
			"start_time":        contest.StartTime,
			"end_time":          contest.EndTime,
			"location":          contest.Location,
			"status":            contest.Status,
			"site_abbreviation": string(contest.SiteAbbreviation), // Assuming it's a string type
		})
	}

	return result, nil
}

// GetContestsOfSpecificSites retrieves contests for specific sites with pagination
func (s *KontestService) GetContestsOfSpecificSites(sites []string, page, perPage int) ([]map[string]string, error) {
	s.fetchHtmlIfNeeded() // Ensure contests are fetched if needed

	// Calculate offset for pagination
	offset := (page - 1) * perPage

	var contests []model.KontestModel

	// Filter contests from the cache based on specific sites
	for _, contest := range s.kontestsCache {
		for _, site := range sites {
			if string(contest.SiteAbbreviation) == site {
				contests = append(contests, contest)
				break
			}
		}
	}

	// Apply pagination
	start := offset
	end := offset + perPage
	if start >= len(contests) {
		return []map[string]string{}, nil // Return an empty result if the offset is beyond available contests
	}
	if end > len(contests) {
		end = len(contests) // Adjust end if it exceeds the slice length
	}

	// Convert to a slice of maps for easier JSON serialization
	result := make([]map[string]string, end-start)
	for i, contest := range contests[start:end] {
		result[i] = map[string]string{
			"id":                contest.ID.String(),
			"name":              contest.Name,
			"url":               contest.URL,
			"start_time":        contest.StartTime,
			"end_time":          contest.EndTime,
			"location":          contest.Location,
			"status":            contest.Status,
			"site_abbreviation": string(contest.SiteAbbreviation),
		}
	}

	return result, nil
}

// Method to check if an update is needed (this should be implemented based on your logic)
func (s *KontestService) shouldUpdate() bool {
	return time.Since(s.lastUpdatedAt) >= updateInterval
}

// Try to acquire the lock for updating
func (s *KontestService) tryUpdate() bool {
	if s.isUpdating.TryLock() {
		return true
	}
	return false
}

func (s *KontestService) fetchHtmlIfNeeded() {
	if s.shouldUpdate() {
		// Try to lock for updating
		if !s.tryUpdate() {
			log.Println("Update is already in progress by another thread.")
			return
		}

		// Perform the fetch operation
		defer s.isUpdating.Unlock() // Ensure that we unlock even if an error occurs
		s.fetchHtml()               // This will perform the fetching and updating logic
	} else {
		log.Println("Update is not required.")
	}
}

func (s *KontestService) PurgeMetadata() {
	database.GetDB().Unscoped().Delete(&model.Metadata{}, "1=1")
}
