package data

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
)

// ContributionEntry struct for storing contributions
type ContributionEntry struct {
	Date string `json:"date"`
	Data int    `json:"data"`
}

// Contributions struct for storing contributions
type Contributions struct {
	Total            int                 `json:"total"`
	ContributionData []ContributionEntry `json:"contributions"`
}

// Request struct for storing contributions
type Request struct {
	Usernames        []string        `json:"usernames"`
	Year             string          `json:"string"`
	Theme            string          `json:"theme"`
	ContributionList []Contributions `json:"contribution_list"`
}

// Aggregate function aggregates contributions of all usernames
func Aggregate(contributionsList []Contributions) (int, []int) {
	aggregateContributions := new(Contributions)
	contributionList := make([]int, 0)

	copier.Copy(&aggregateContributions, &contributionsList[0])

	for i := 1; i < len(contributionsList); i++ {
		aggregateContributions.Total += contributionsList[i].Total

		for j, contributionEntry := range contributionsList[i].ContributionData {
			aggregateContributions.ContributionData[j].Data += contributionEntry.Data
		}
	}

	for _, contribution := range aggregateContributions.ContributionData {
		contributionList = append(contributionList, contribution.Data)
	}

	return aggregateContributions.Total, contributionList
}

// GetRawPage function fetches the raw HTML of GitHub user's page
func GetRawPage(username string, year string) (string, error) {
	url := "https://www.github.com/users/" + username + "/contributions?from=" + year + "-01-01&to=" + year + "-12-31"
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {
		return "", fmt.Errorf("user not found")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	rawPage := string(body)
	return rawPage, nil
}

// ParseContributionsData function parses the required contributions data
func ParseContributionsData(rawHTML string, year string) Contributions {
	var contributions Contributions

	contributions.Total = parseTotal(rawHTML)
	dateData := parseDateData(rawHTML, year)

	contributions.ContributionData = dateData

	return contributions
}

// parseTotal function parses the total contributions
func parseTotal(rawHTML string) int {
	r, _ := regexp.Compile("[0-9]+ contributions")
	totalContributionsRaw := r.FindString(rawHTML)
	if totalContributionsRaw == "" {
		return 0
	}

	totalContributionsRaw = strings.Fields(totalContributionsRaw)[0]
	totalContributions, _ := strconv.Atoi(totalContributionsRaw)

	return totalContributions
}

// parseDateData function parses the contributions date-data
func parseDateData(rawHTML string, year string) []ContributionEntry {
	regexString := "data-count=\"[0-9]{1,3}\" data-date=\"" + year + "-[0-9]{2}-[0-9]{2}\""

	r := regexp.MustCompile(regexString)
	allDatesContributions := r.FindAllString(rawHTML, -1)

	contributionDateData := make([]ContributionEntry, 0)

	for _, singleDateContribution := range allDatesContributions {
		parts := strings.Split(singleDateContribution, "\"")

		contributionDate := parts[3]                  // Extracts date
		contributionData, _ := strconv.Atoi(parts[1]) // Extracts contribution for the date

		contributionDateData = append(contributionDateData, ContributionEntry{contributionDate, contributionData})
	}
	return contributionDateData
}
