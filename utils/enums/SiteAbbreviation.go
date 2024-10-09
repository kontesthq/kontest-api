package enums

type SiteAbbreviation string

// Define the abbreviations as constants
const (
	AtCoder       SiteAbbreviation = "AtCoder"
	CodeChef      SiteAbbreviation = "CodeChef"
	CodeForces    SiteAbbreviation = "CodeForces"
	CodeForcesGym SiteAbbreviation = "CodeForcesGym"
	CodingNinjas  SiteAbbreviation = "CodingNinjas"
	CSAcademy     SiteAbbreviation = "CSAcademy"
	GeeksForGeeks SiteAbbreviation = "GeeksForGeeks"
	HackerEarth   SiteAbbreviation = "HackerEarth"
	HackerRank    SiteAbbreviation = "HackerRank"
	LeetCode      SiteAbbreviation = "LeetCode"
	ProjectEuler  SiteAbbreviation = "ProjectEuler"
	TopCoder      SiteAbbreviation = "TopCoder"
	Toph          SiteAbbreviation = "Toph"
	YukiCoder     SiteAbbreviation = "YukiCoder"
	CupsOnline    SiteAbbreviation = "CupsOnline"
	RoboContest   SiteAbbreviation = "RoboContest"
	CTFtime       SiteAbbreviation = "CTFtime"
	LightOJ       SiteAbbreviation = "LightOJ"
	UCup          SiteAbbreviation = "UCup"
	Kaggle        SiteAbbreviation = "Kaggle"
	DMOJ          SiteAbbreviation = "DMOJ"
	TLX           SiteAbbreviation = "TLX"
	CodeRun       SiteAbbreviation = "CodeRun"
	Eolymp        SiteAbbreviation = "eOlymp"
	ICPCGlobal    SiteAbbreviation = "ICPCGlobal"
	Luogu         SiteAbbreviation = "Luogu"
	SPOJ          SiteAbbreviation = "SPOJ"
	GSU           SiteAbbreviation = "GSU"
)

// GetAbbreviation returns the site abbreviation based on the input location
func GetAbbreviation(location string) SiteAbbreviation {
	switch location {
	case "atcoder.jp":
		return AtCoder
	case "codechef.com":
		return CodeChef
	case "codeforces.com":
		return CodeForces
	case "codingninjas.com", "codingninjas.com/codestudio":
		return CodingNinjas
	case "csacademy.com":
		return CSAcademy
	case "geeksforgeeks.org":
		return GeeksForGeeks
	case "hackerearth.com":
		return HackerEarth
	case "hackerrank.com":
		return HackerRank
	case "leetcode.com":
		return LeetCode
	case "projecteuler.net":
		return ProjectEuler
	case "topcoder.com":
		return TopCoder
	case "toph.com", "toph.co":
		return Toph
	case "yukicoder.me":
		return YukiCoder
	case "cups.online":
		return CupsOnline
	case "robocontest.uz":
		return RoboContest
	case "ctftime.org":
		return CTFtime
	case "lightoj.com":
		return LightOJ
	case "ucup.ac":
		return UCup
	case "kaggle.com":
		return Kaggle
	case "dmoj.ca":
		return DMOJ
	case "tlx.toki.id":
		return TLX
	case "coderun.yandex.ru":
		return CodeRun
	case "eolymp.com":
		return Eolymp
	case "icpc.global":
		return ICPCGlobal
	case "luogu.com.cn":
		return Luogu
	case "spoj.com":
		return SPOJ
	case "dl.gsu.by":
		return GSU
	default:
		return SiteAbbreviation(location) // return the input location if no match
	}
}
