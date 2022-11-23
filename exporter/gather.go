package exporter

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// gatherData - Collects the data from the API and stores into struct
func (e *Exporter) gatherData() ([]*Datum, error) {

	data := []*Datum{}

	responses, err := asyncHTTPGets(e.TargetURLs(), e.APIToken())

	if err != nil {
		return data, err
	}

	for _, response := range responses {

		// Github can at times present an array, or an object for the same data set.
		// This code checks handles this variation.
		if isArray(response.body) {
			ds := []*Datum{}
			json.Unmarshal(response.body, &ds)
			data = append(data, ds...)
		} else {
			d := new(Datum)

			if strings.Contains(response.url, "/repos/") {
				getReleases(e, response.url, &d.Releases)
				getPRs(e, response.url, &d.Pulls)
				getClones(e, response.url, &d.Clones)
				getReferralPaths(e, response.url, &d.ReferralPaths)
				getReferralSources(e, response.url, &d.ReferralSources)
				getPageViews(e, response.url, &d.PageViews)
			}
			json.Unmarshal(response.body, &d)
			data = append(data, d)
		}

		log.Infof("API data fetched for repository: %s", response.url)
	}

	//return data, rates, err
	return data, nil

}

// getRates obtains the rate limit data for requests against the github API.
// Especially useful when operating without oauth and the subsequent lower cap.
func (e *Exporter) getRates() (*RateLimits, error) {
	u := *e.APIURL()
	u.Path = path.Join(u.Path, "rate_limit")

	resp, err := getHTTPResponse(u.String(), e.APIToken())
	if err != nil {
		return &RateLimits{}, err
	}
	defer resp.Body.Close()

	// Triggers if rate-limiting isn't enabled on private Github Enterprise installations
	if resp.StatusCode == 404 {
		return &RateLimits{}, fmt.Errorf("Rate Limiting not enabled in GitHub API")
	}

	limit, err := strconv.ParseFloat(resp.Header.Get("X-RateLimit-Limit"), 64)

	if err != nil {
		return &RateLimits{}, err
	}

	rem, err := strconv.ParseFloat(resp.Header.Get("X-RateLimit-Remaining"), 64)

	if err != nil {
		return &RateLimits{}, err
	}

	reset, err := strconv.ParseFloat(resp.Header.Get("X-RateLimit-Reset"), 64)

	if err != nil {
		return &RateLimits{}, err
	}

	return &RateLimits{
		Limit:     limit,
		Remaining: rem,
		Reset:     reset,
	}, err

}

func getReleases(e *Exporter, url string, data *[]Release) {
	i := strings.Index(url, "?")
	baseURL := url[:i]
	releasesURL := baseURL + "/releases"
	releasesResponse, err := asyncHTTPGets([]string{releasesURL}, e.APIToken())

	if err != nil {
		log.Errorf("Unable to obtain releases from API, Error: %s", err)
	}

	for _, r := range releasesResponse {
		d := []Release{}
		json.Unmarshal(r.body, &d)
		*data = append(*data, d...)
	}
}

func getPRs(e *Exporter, url string, data *[]Pull) {
	i := strings.Index(url, "?")
	baseURL := url[:i]
	pullsURL := baseURL + "/pulls"
	pullsResponse, err := asyncHTTPGets([]string{pullsURL}, e.APIToken())

	if err != nil {
		log.Errorf("Unable to obtain pull requests from API, Error: %s", err)
	}

	for _, r := range pullsResponse {
		d := []Pull{}
		json.Unmarshal(r.body, &d)
		*data = append(*data, d...)
	}
}

func getClones(e *Exporter, url string, data *Clones) {
	i := strings.Index(url, "?")
	baseURL := url[:i]
	clonesURL := baseURL + "/traffic/clones"
	clonesResponse, err := asyncHTTPGets([]string{clonesURL}, e.APIToken())

	if err != nil {
		log.Errorf("Unable to obtain repository clones from API, Error: %s", err)
	}

	json.Unmarshal(clonesResponse[0].body, &data)
}

func getReferralPaths(e *Exporter, url string, data *[]ReferralPath) {
	i := strings.Index(url, "?")
	baseURL := url[:i]
	referralPathsURL := baseURL + "/traffic/popular/paths"
	referralPathsResponse, err := asyncHTTPGets([]string{referralPathsURL}, e.APIToken())

	if err != nil {
		log.Errorf("Unable to obtain referral paths from API, Error: %s", err)
	}

	json.Unmarshal(referralPathsResponse[0].body, &data)
}

func getReferralSources(e *Exporter, url string, data *[]ReferralSource) {
	i := strings.Index(url, "?")
	baseURL := url[:i]
	referralSourceURL := baseURL + "/traffic/popular/referrers"
	referralSourceResponse, err := asyncHTTPGets([]string{referralSourceURL}, e.APIToken())

	if err != nil {
		log.Errorf("Unable to obtain referral sources from API, Error: %s", err)
	}

	json.Unmarshal(referralSourceResponse[0].body, &data)
}

func getPageViews(e *Exporter, url string, data *PageViews) {
	i := strings.Index(url, "?")
	baseURL := url[:i]
	pageViewsURL := baseURL + "/traffic/views"
	pageViewsResponse, err := asyncHTTPGets([]string{pageViewsURL}, e.APIToken())

	if err != nil {
		log.Errorf("Unable to obtain pull requests from API, Error: %s", err)
	}

	json.Unmarshal(pageViewsResponse[0].body, &data)
}

// isArray simply looks for key details that determine if the JSON response is an array or not.
func isArray(body []byte) bool {

	isArray := false

	for _, c := range body {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		isArray = c == '['
		break
	}

	return isArray

}
