package utils

import "kontest-api/utils/enums"

// SiteUtils provides utility functions related to sites
type SiteUtils struct{}

// GetSiteAbbreviationFromSite returns the abbreviation for a given site location
func (su *SiteUtils) GetSiteAbbreviationFromSite(location string) string {
	abbreviation := enums.GetAbbreviation(location)
	return string(abbreviation)
}
