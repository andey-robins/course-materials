package scrape

import (
	"regexp"
)

//==========================================================================\\
// || GLOBAL DATA STRUCTURES  ||

type FileInfo struct {
	Filename string `json:"filename"`
	Location string `json:"location"`
}

var Files []FileInfo

var regexes = []*regexp.Regexp{
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?)i.txt`),
}

// END GLOBAL VARIABLES
//==========================================================================//

//==========================================================================\\
// || HELPER FUNCTIONS TO MANIPULATE THE REGULAR EXPRESSIONS ||

func resetRegEx() {
	regexes = []*regexp.Regexp{
		regexp.MustCompile(`(?i)password`),
		regexp.MustCompile(`(?i)kdb`),
		regexp.MustCompile(`(?i)login`),
	}
}

func clearRegEx() {
	regexes = make([]*regexp.Regexp, 0)
}

func addRegEx(regex string) {
	regexes = append(regexes, regexp.MustCompile(regex))
}

//==========================================================================//
