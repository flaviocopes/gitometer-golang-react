package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/flaviocopes/gitometer/server/common"
	"github.com/flaviocopes/gitometer/server/db"
	gogithub "github.com/google/go-github/github"
	"github.com/jinzhu/now"
	"golang.org/x/oauth2"
)

var clientV3 *gogithub.Client

func getClientV3() *gogithub.Client {
	if clientV3 == nil {
		ctx := context.Background()
		at := os.Getenv("GITOMETER_GITHUB_ACCESS_TOKEN")
		if at == "" {
			panic("You need to set the GitHub access token as GITOMETER_GITHUB_ACCESS_TOKEN environment variable")
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: at},
		)
		tc := oauth2.NewClient(ctx, ts)
		clientV3 = gogithub.NewClient(tc)
	}

	return clientV3
}

func getBasicRepoInfo(owner, name string) *common.Repository {
	ctx := context.Background()
	repo, _, err := getClientV3().Repositories.Get(ctx, owner, name)
	if err != nil {
		panic(err)
	}

	r := common.Repository{}
	r.ID = *repo.ID
	r.Name = *repo.Name
	r.OwnerName = owner
	r.DefaultBranch = *repo.DefaultBranch
	r.CreatedAt = (*repo.CreatedAt).Format(time.RFC3339)
	r.Initialized = false
	r.Description = *repo.Description
	r.RepoAge = monthsCountSince((*repo.CreatedAt).Time)
	r.TotalStars = *repo.StargazersCount
	r.TotalCommits = getTotalCommits(owner, name)
	r.CommitsCountLast12Months, r.CommitsCountLast4Weeks, r.CommitsCountLastWeek = getCommitsData(owner, name)
	r.StarsCountLast12Months, r.StarsCountLast4Weeks, r.StarsCountLastWeek, r.StarsPerMonth = getStarsData(owner, name)

	return &r
}

func getCommitsData(owner, name string) (int, int, int) {
	data, _, err := getClientV3().Repositories.ListParticipation(context.Background(), owner, name)
	if err != nil {
		log.Fatalf("Repositories.ListParticipation returned error: %v", err)
	}
	w := data.All
	commitsCountLast12Months := 0
	for _, v := range w {
		commitsCountLast12Months += v
	}
	commitsCountLast4Weeks := 0
	for _, v := range w[len(w)-4:] {
		commitsCountLast4Weeks += v
	}
	commitsCountLastWeek := w[len(w)-1]
	return commitsCountLast12Months, commitsCountLast4Weeks, commitsCountLastWeek
}

func getTotalCommits(owner, name string) int {
	opt := &gogithub.CommitsListOptions{
		ListOptions: gogithub.ListOptions{PerPage: 30},
	}
	_, resp, err := getClientV3().Repositories.ListCommits(context.Background(), owner, name, opt)
	if err != nil {
		log.Fatalf("Repositories.ListCommits returned error: %v", err)
	}
	pagesCount := resp.LastPage
	tot := (pagesCount - 1) * 30 // 30 items per page until the last one
	opt.Page = resp.LastPage
	commits, _, err := getClientV3().Repositories.ListCommits(context.Background(), owner, name, opt)
	if err != nil {
		log.Fatalf("Repositories.ListCommits returned error: %v", err)
	}
	tot = tot + len(commits) //the last page contains the ramaining few commits, <= 30
	return tot
}

// monthsCountSince calculates the months between now
// and the createdAtTime time.Time value passed
func monthsCountSince(createdAtTime time.Time) int {
	now := time.Now()
	months := 0
	month := createdAtTime.Month()
	for createdAtTime.Before(now) {
		createdAtTime = createdAtTime.Add(time.Hour * 24)
		nextMonth := createdAtTime.Month()
		if nextMonth != month {
			months++
		}
		month = nextMonth
	}

	return months
}

// AddRepoToDb adds a repository to the database. It's the only
// exported method of this package.
func AddRepoToDb(owner, name string) {
	repo := getBasicRepoInfo(owner, name)
	err := db.AddNewRepo(owner, name, repo)
	if err != nil {
		panic(err)
	}
}

type yearmonth struct{ Year, Month int }
type yearweek struct{ Year, Week int }

func getStarsData(owner, name string) (int, int, int, string) {

	dateTimeNow := time.Now()
	dateTimeLastWeek := dateTimeNow.AddDate(0, 0, -7)

	startWeek := now.New(dateTimeLastWeek).BeginningOfWeek()
	endWeek := now.New(dateTimeLastWeek).EndOfWeek()

	lastSundayDate := time.Time(endWeek)
	dateStartLastWeek := time.Time(startWeek)

	dateStartLast4Weeks := startWeek.AddDate(0, 0, -7*3)
	dateStartLast12Months := startWeek.AddDate(0, 0, -7*51)

	opt := gogithub.ListOptions{PerPage: 10}
	stargazers, resp, err := clientV3.Activity.ListStargazers(context.Background(), owner, name, &opt)
	if err != nil {
		panic(err)
	}
	var stars []time.Time
	for _, v := range stargazers {
		stars = append(stars, v.StarredAt.Time)
	}

	stars = reverseStarsSlice(stars)
	results := stars

	continueFlag := true

	if resp.LastPage != resp.FirstPage { //more than one page
		opt.Page = resp.LastPage
		stargazers, resp, err := clientV3.Activity.ListStargazers(context.Background(), owner, name, &opt)
		if err != nil {
			panic(err)
		}
		for _, v := range stargazers {
			stars = append(stars, v.StarredAt.Time)
		}

		stars = reverseStarsSlice(stars)
		results = append(results, stars...)
		prevPage := resp.PrevPage

		for {
			if !continueFlag {
				break
			}
			if prevPage == 0 {
				break
			}

			opt.Page = prevPage
			fmt.Println(prevPage)

			stargazers, resp, err := clientV3.Activity.ListStargazers(context.Background(), owner, name, &opt)
			if err != nil {
				panic(err)
			}

			var s []time.Time
			for _, v := range stargazers {
				s = append(s, v.StarredAt.Time)
			}

			s = reverseStarsSlice(s)
			results = append(results, s...)
			time.Sleep(time.Second * 1)

			prevPage = resp.PrevPage
		}

	}

	starsCountLastWeek := 0
	starsCountLast4Weeks := 0
	starsCountLast12Months := 0

	weeklyData := make(map[yearweek]int)
	starringData := make(map[yearmonth]int)
	month := 0
	year := 0

	for _, result := range results {
		starredAt := result
		month = int(result.Month())
		year = int(result.Year())

		if _, ok := starringData[yearmonth{year, month}]; ok {
			starringData[yearmonth{year, month}]++
		} else {
			starringData[yearmonth{year, month}] = 1
		}

		if starredAt.After(lastSundayDate) {
			// drop newer stars
			continue
		}

		if starredAt.After(dateStartLastWeek) {
			starsCountLastWeek++
		}
		if starredAt.After(dateStartLast4Weeks) {
			starsCountLast4Weeks++
		}
		if starredAt.After(dateStartLast12Months) {
			starsCountLast12Months++
		}

		_, week := starredAt.ISOWeek()

		if _, ok := weeklyData[yearweek{year, week}]; ok {
			weeklyData[yearweek{year, week}]++
		} else {
			weeklyData[yearweek{year, week}] = 1
		}
	}

	starsPerMonth := prepareDataForGraph(starringData)

	return starsCountLast12Months, starsCountLast4Weeks, starsCountLastWeek, starsPerMonth
}

func fillMissingMonths(data map[yearmonth]int) (map[yearmonth]int, int, int, int, int) {
	keys := make([]yearmonth, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var lowestYear, highestYear, lowestMonthInLowestYear, highestMonthInHighestYear int

	for _, k := range keys {
		if highestYear == 0 || k.Year >= highestYear {
			highestYear = k.Year
			if highestMonthInHighestYear == 0 || k.Month > highestMonthInHighestYear {
				highestMonthInHighestYear = k.Month
			}
		}
		if lowestYear == 0 || k.Year <= lowestYear {
			lowestYear = k.Year
			if lowestMonthInLowestYear == 0 || k.Month < lowestMonthInLowestYear {
				lowestMonthInLowestYear = k.Month
			}
		}
	}

	for month := lowestMonthInLowestYear; month <= 12; month++ {
		if _, ok := data[yearmonth{lowestYear, month}]; !ok {
			data[yearmonth{lowestYear, month}] = 0
		}
	}

	for year := lowestYear + 1; year < highestYear; year++ {
		for month := 1; month <= 12; month++ {
			if _, ok := data[yearmonth{year, month}]; !ok {
				data[yearmonth{year, month}] = 0
			}
		}
	}

	for month := highestMonthInHighestYear; month > 0; month-- {
		if _, ok := data[yearmonth{highestYear, month}]; !ok {
			data[yearmonth{highestYear, month}] = 0
		}
	}

	return data, lowestYear, highestYear, lowestMonthInLowestYear, highestMonthInHighestYear
}

func generateLabelsForGraph(data map[yearmonth]int, lowestYear, highestYear, lowestMonthInLowestYear, highestMonthInHighestYear int) []string {
	var labels []string

	for month := lowestMonthInLowestYear; month <= 12; month++ {
		if _, ok := data[yearmonth{lowestYear, month}]; ok {
			labels = append(labels, fmt.Sprintf("%d %d", month, lowestYear))
		}
	}

	for year := lowestYear + 1; year < highestYear; year++ {
		for month := 1; month <= 12; month++ {
			if _, ok := data[yearmonth{year, month}]; ok {
				labels = append(labels, fmt.Sprintf("%d %d", month, year))
			}
		}
	}

	for month := 1; month <= highestMonthInHighestYear; month++ {
		if _, ok := data[yearmonth{highestYear, month}]; ok {
			labels = append(labels, fmt.Sprintf("%d %d", month, highestYear))
		}
	}

	return labels
}

func generateDataForGraph(data map[yearmonth]int, lowestYear, highestYear, lowestMonthInLowestYear, highestMonthInHighestYear int) []int {
	var commitsCount []int
	total := 0

	year := lowestYear
	for month := lowestMonthInLowestYear; month <= 12; month++ {
		total += data[yearmonth{lowestYear, month}]
		commitsCount = append(commitsCount, total)
	}

	for year = lowestYear + 1; year < highestYear; year++ {
		for month := 1; month <= 12; month++ {
			total += data[yearmonth{year, month}]
			commitsCount = append(commitsCount, total)
		}
	}

	year = highestYear
	for month := highestMonthInHighestYear; month > 0; month-- {
		total += data[yearmonth{year, month}]
		commitsCount = append(commitsCount, total)
	}

	return commitsCount
}

func prepareDataForGraph(data map[yearmonth]int) string {
	preparedData, lowestYear, highestYear, lowestMonthInLowestYear, highestMonthInHighestYear := fillMissingMonths(data)

	type dataForGraph struct {
		Labels []string `json:"labels"`
		Data   []int    `json:"data"`
	}

	graphLabels := generateLabelsForGraph(preparedData, lowestYear, highestYear, lowestMonthInLowestYear, highestMonthInHighestYear)
	graphData := generateDataForGraph(preparedData, lowestYear, highestYear, lowestMonthInLowestYear, highestMonthInHighestYear)

	c, err := json.Marshal(dataForGraph{graphLabels, graphData})
	if err != nil {
		panic(err)
	}

	return string(c)
}

func reverseStarsSlice(s []time.Time) []time.Time {
	var n []time.Time
	for i := len(s) - 1; i >= 0; i-- {
		n = append(n, s[i])
	}
	return n
}
