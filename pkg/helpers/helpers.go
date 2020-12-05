package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ContributionEntry struct for storing contributions
type ContributionEntry struct {
	Date string
	Data int
}

// Contributions struct for storing contributions
type Contributions struct {
	Total            int
	ContributionData []ContributionEntry
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
	dateData, err := parseContributionDataDate(rawHTML)
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

	totalContributionsRaw = strings.Fields(totalContributionsRaw)[0]
	totalContributions, _ := strconv.Atoi(totalContributionsRaw)

	return totalContributions
}

// ParseContributionsData function parses the contributions date-data
func parseContributionDataDate(rawHTML string) ([]ContributionEntry, error) {
	r, _ := regexp.Compile("data-count=\"[0-9]{1,3}\" data-date=\"2020-[0-9]{2}-[0-9]{2}\"")
	allContributionsRaw := r.FindAllString(rawHTML, -1)

	var contributionDataDate []ContributionEntry

	for _, singleDate := range allContributionsRaw {
		parts := strings.Split(singleDate, "\"")

		date := parts[3]
		contributions, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Println("Error in converting string to int")
			return nil, err
		}

		contributionDataDate = append(contributionDataDate, ContributionEntry{date, contributions})
	}

	return contributionDataDate, nil
}
