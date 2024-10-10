package enums

// SiteAbbreviation represents a site abbreviation with its corresponding URL.
type SiteAbbreviation struct {
	Name string
	URL  string
}

// List of all site abbreviations.
var abbreviations = []SiteAbbreviation{
	{"AtCoder", "atcoder.jp"},
	{"CodeChef", "codechef.com"},
	{"CodeForces", "codeforces.com"},
	{"CodeForcesGym", "codeforces.com/gym"},
	{"CodingNinjas", "codingninjas.com"},
	{"CSAcademy", "csacademy.com"},
	{"GeeksForGeeks", "geeksforgeeks.org"},
	{"HackerEarth", "hackerearth.com"},
	{"HackerRank", "hackerrank.com"},
	{"LeetCode", "leetcode.com"},
	{"ProjectEuler", "projecteuler.net"},
	{"TopCoder", "topcoder.com"},
	{"Toph", "toph.com"},
	{"YukiCoder", "yukicoder.me"},
	{"CupsOnline", "cups.online"},
	{"RoboContest", "robocontest.uz"},
	{"CTFtime", "ctftime.org"},
	{"LightOJ", "lightoj.com"},
	{"UCup", "ucup.ac"},
	{"Kaggle", "kaggle.com"},
	{"DMOJ", "dmoj.ca"},
	{"TLX", "tlx.toki.id"},
	{"CodeRun", "coderun.yandex.ru"},
	{"Eolymp", "eolymp.com"},
	{"ICPCGlobal", "icpc.global"},
	{"Luogu", "luogu.com.cn"},
	{"SPOJ", "spoj.com"},
	{"GSU", "dl.gsu.by"},
}

// GetAbbreviation returns the site abbreviation based on the input location.
func GetAbbreviation(location string) string {
	for _, abbr := range abbreviations {
		if abbr.URL == location {
			return abbr.Name
		}
	}
	return location // Return the input location if no match.
}

// GetAllAbbreviations returns a list of all abbreviations.
func GetAllAbbreviations() []string {
	var names []string
	for _, abbr := range abbreviations {
		names = append(names, abbr.Name)
	}
	return names
}
