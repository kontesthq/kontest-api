package controllers

import (
	"encoding/json"
	"kontest-api/utils"
	"kontest-api/utils/enums"
	"net/http"
	"strconv"
	"strings"
)

func GetAllKontests(w http.ResponseWriter, r *http.Request) {
	kontestService := utils.GetDependencies().KontestService

	// Parse query parameters
	rawSites := r.URL.Query().Get("sites") // Get the sites parameter as a single string
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	var siteList []string
	if rawSites != "" {
		// Split the comma-separated sites into a slice
		siteList = strings.Split(rawSites, ",")
	}

	// Convert page and perPage to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	perPage := 10 // Default perPage
	if perPageStr != "" {
		perPage, err = strconv.Atoi(perPageStr)
		if err != nil || perPage <= 0 {
			http.Error(w, "Invalid per_page number", http.StatusBadRequest)
			return
		}
	}

	var contests []map[string]string
	if len(siteList) == 0 {
		// If there are no specific sites, get all contests
		contests, err = kontestService.GetContests(page, perPage)
	} else {
		// If there are specific sites, fetch contests for those sites
		contests, err = kontestService.GetContestsOfSpecificSites(siteList, page, perPage)
	}

	if err != nil {
		http.Error(w, "Failed to get contests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the contests in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contests)
}

func PurgeMetadata(w http.ResponseWriter, r *http.Request) {
	kontestService := utils.GetDependencies().KontestService

	kontestService.PurgeMetadata()

	json.NewEncoder(w).Encode(map[string]string{"message": "Metadata purged successfully"})
}

func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Service is healthy"))
}

func GetSupportedSites(writer http.ResponseWriter, request *http.Request) {
	supportedSites := enums.GetAllAbbreviations()
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(supportedSites)
}
