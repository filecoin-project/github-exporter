package exporter

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// AddMetrics - Add's all of the metrics to a map of strings, returns the map.
func AddMetrics() map[string]*prometheus.Desc {

	APIMetrics := make(map[string]*prometheus.Desc)

	APIMetrics["Stars"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "stars"),
		"Total number of Stars for given repository",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["OpenIssues"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "open_issues"),
		"Total number of open issues for given repository",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["PullRequestCount"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "pull_request_count"),
		"Total number of pull requests for given repository",
		[]string{"repo"}, nil,
	)
	APIMetrics["Watchers"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "watchers"),
		"Total number of watchers/subscribers for given repository",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["Forks"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "forks"),
		"Total number of forks for given repository",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["Size"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "size_kb"),
		"Size in KB for given repository",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["RepoClonesDaily"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "clones_daily"),
		"Get the total number of clones for today",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["RepoClonesBiweekly"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "clones_biweekly"),
		"Get the total number of clones for the last 14 days",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["UniqueRepoClonesDaily"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "unique_clones_daily"),
		"Get the total number of unique clones for today",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["UniqueRepoClonesBiweekly"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "unique_clones_biweekly"),
		"Get the total number of unique clones for the last 14 days",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["PageViewsDaily"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "page_views_daily"),
		"Get the total number of views for today",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["PageViewsBiweekly"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "page_views_biweekly"),
		"Get the total number of views for the last 14 days",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["UniquePageViewsDaily"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "unique_page_views_daily"),
		"Get the total number of unique views for today",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["UniquePageViewsBiweekly"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "unique_page_views_biweekly"),
		"Get the total number of unique views for the last 14 days",
		[]string{"repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["ReferralSources"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "referral_sources_biweekly"),
		"Get the total visitor count top 10 referrers over the last 14 days",
		[]string{"source", "repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["UniqueReferralSources"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "unique_referral_sources_biweekly"),
		"Get the unique visitor count for top 10 referrers over the last 14 days",
		[]string{"source", "repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["ReferralPaths"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "referral_paths_biweekly"),
		"Get the total visitor count top 10 popular contents over the last 14 days",
		[]string{"path", "repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["UniqueReferralPaths"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "unique_referral_paths_biweekly"),
		"Get the unique visitor count for the top 10 popular contents over the last 14 days",
		[]string{"path", "repo", "user", "private", "fork", "archived", "license", "language"}, nil,
	)
	APIMetrics["ReleaseDownloads"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "repo", "release_downloads"),
		"Download count for a given release",
		[]string{"repo", "user", "release", "name", "created_at"}, nil,
	)
	APIMetrics["Limit"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "rate", "limit"),
		"Number of API queries allowed in a 60 minute window",
		[]string{}, nil,
	)
	APIMetrics["Remaining"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "rate", "remaining"),
		"Number of API queries remaining in the current window",
		[]string{}, nil,
	)
	APIMetrics["Reset"] = prometheus.NewDesc(
		prometheus.BuildFQName("github", "rate", "reset"),
		"The time at which the current rate limit window resets in UTC epoch seconds",
		[]string{}, nil,
	)

	return APIMetrics
}

func (e *Exporter) processMetrics(data []*Datum, rates *RateLimits, ch chan<- prometheus.Metric) error {

	for _, x := range data {
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["Stars"], prometheus.GaugeValue, x.Stars, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["Forks"], prometheus.GaugeValue, x.Forks, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["Watchers"], prometheus.GaugeValue, x.Watchers, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["Size"], prometheus.GaugeValue, x.Size, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)

		for _, release := range x.Releases {
			for _, asset := range release.Assets {
				ch <- prometheus.MustNewConstMetric(e.APIMetrics["ReleaseDownloads"], prometheus.GaugeValue, float64(asset.Downloads), x.Name, x.Owner.Login, release.Name, asset.Name, asset.CreatedAt)
			}
		}

		prCount := float64(len(x.Pulls))
		// The xOpenIssues field in the Github API includes all Open Issues and all Open PRs, so subtract the Open PRs to get only the Open Issue count.
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["OpenIssues"], prometheus.GaugeValue, (x.OpenIssues - prCount), x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["PullRequestCount"], prometheus.GaugeValue, prCount, x.Name)

		ch <- prometheus.MustNewConstMetric(e.APIMetrics["RepoClonesBiweekly"], prometheus.GaugeValue, x.Clones.Count, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["UniqueRepoClonesBiweekly"], prometheus.GaugeValue, x.Clones.Uniques, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		for _, dailyClones := range x.Clones.Clones {
			if time.Now().Format(time.RFC3339) == dailyClones.Timestamp {
				ch <- prometheus.MustNewConstMetric(e.APIMetrics["RepoClonesDaily"], prometheus.GaugeValue, dailyClones.Count, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
				ch <- prometheus.MustNewConstMetric(e.APIMetrics["UniqueRepoClonesDaily"], prometheus.GaugeValue, dailyClones.Uniques, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
			}
		}

		ch <- prometheus.MustNewConstMetric(e.APIMetrics["PageViewsBiweekly"], prometheus.GaugeValue, x.PageViews.Count, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		ch <- prometheus.MustNewConstMetric(e.APIMetrics["UniquePageViewsBiweekly"], prometheus.GaugeValue, x.PageViews.Uniques, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		for _, dailyPageViews := range x.PageViews.Views {
			if time.Now().Format(time.RFC3339) == dailyPageViews.Timestamp {
				ch <- prometheus.MustNewConstMetric(e.APIMetrics["PageViewsDaily"], prometheus.GaugeValue, dailyPageViews.Count, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
				ch <- prometheus.MustNewConstMetric(e.APIMetrics["UniquePageViewsDaily"], prometheus.GaugeValue, dailyPageViews.Uniques, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
			}
		}

		for _, source := range x.ReferralSources {
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["ReferralSources"], prometheus.GaugeValue, source.Count, source.Referrer, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["UniqueReferralSources"], prometheus.GaugeValue, source.Uniques, source.Referrer, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		}

		for _, path := range x.ReferralPaths {
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["ReferralPaths"], prometheus.GaugeValue, path.Count, path.Path, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
			ch <- prometheus.MustNewConstMetric(e.APIMetrics["UniqueReferralPaths"], prometheus.GaugeValue, path.Uniques, path.Path, x.Name, x.Owner.Login, strconv.FormatBool(x.Private), strconv.FormatBool(x.Fork), strconv.FormatBool(x.Archived), x.License.Key, x.Language)
		}
	}

	ch <- prometheus.MustNewConstMetric(e.APIMetrics["Limit"], prometheus.GaugeValue, rates.Limit)
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["Remaining"], prometheus.GaugeValue, rates.Remaining)
	ch <- prometheus.MustNewConstMetric(e.APIMetrics["Reset"], prometheus.GaugeValue, rates.Reset)

	return nil
}
