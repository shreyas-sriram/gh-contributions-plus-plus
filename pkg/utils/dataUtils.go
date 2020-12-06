package utils

import (
	"fmt"
	"io/ioutil"
	"log"
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

var (
	// ContributionList contains the contributions of the year as an integer array
	ContributionList []int
)

// AggregateContributions function aggregates contributions of all usernames
func AggregateContributions(contributionsList []Contributions) ([]int, error) {
	var aggregateContributions Contributions

	copier.Copy(&aggregateContributions, &contributionsList[0])

	for i := 1; i < len(contributionsList); i++ {
		aggregateContributions.Total += contributionsList[i].Total

		for j, contributionEntry := range contributionsList[i].ContributionData {
			aggregateContributions.ContributionData[j].Data += contributionEntry.Data
		}
	}

	for _, contribution := range aggregateContributions.ContributionData {
		ContributionList = append(ContributionList, contribution.Data)
	}

	return ContributionList, nil
}

// GetRawPage function fetches the raw HTML of GitHub user's page
func GetRawPage(username string) (string, error) {
	// TODO support arbitrary year
	res, err := http.Get("https://www.github.com/users/" + username + "/contributions?from=2020-01-01&to=2020-12-31")
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
func ParseContributionsData(rawHTML string) (Contributions, error) {
	var contributions Contributions

	contributions.Total = parseTotalContributions(rawHTML)
	dateData, err := parseContributionDateData(rawHTML)
	if err != nil {
		log.Println("Error in converting string to int")
	}

	contributions.ContributionData = dateData

	return contributions, nil
}

// ParseContributionsData function parses the total contributions
func parseTotalContributions(rawHTML string) int {
	r, _ := regexp.Compile("[0-9]+ contributions")
	totalContributionsRaw := r.FindString(rawHTML)
	if totalContributionsRaw == "" {
		return 0
	}

	totalContributionsRaw = strings.Fields(totalContributionsRaw)[0]
	totalContributions, _ := strconv.Atoi(totalContributionsRaw)

	return totalContributions
}

// ParseContributionsData function parses the contributions date-data
func parseContributionDateData(rawHTML string) ([]ContributionEntry, error) {
	r, _ := regexp.Compile("data-count=\"[0-9]{1,3}\" data-date=\"2020-[0-9]{2}-[0-9]{2}\"")
	allDatesContributions := r.FindAllString(rawHTML, -1)

	contributionDateData := make([]ContributionEntry, 0)

	for _, singleDateContribution := range allDatesContributions {
		parts := strings.Split(singleDateContribution, "\"")

		contributionDate := parts[3]                    // Extracts date
		contributionData, err := strconv.Atoi(parts[1]) // Extracts contribution for the date
		if err != nil {
			log.Println("Error in converting string to int")
			return nil, err
		}

		contributionDateData = append(contributionDateData, ContributionEntry{contributionDate, contributionData})
	}

	return contributionDateData, nil
}
