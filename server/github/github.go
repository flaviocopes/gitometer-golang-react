package main

import (
	"context"
	"os"

	"github.com/flaviocopes/gitometer/server/db"
	gogithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	owner := "flaviocopes"
	name := "gitometer-laravel"
	addRepoToDb(owner, name)
}

var client *gogithub.Client

func getClient() *gogithub.Client {
	if client == nil {
		ctx := context.Background()
		at := os.Getenv("GITOMETER_GITHUB_ACCESS_TOKEN")
		if at == "" {
			panic("You need to set the GitHub access token as GITOMETER_GITHUB_ACCESS_TOKEN environment variable")
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: at},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = gogithub.NewClient(tc)
	}

	return client
}

func getBasicRepoInfo(owner, name string) *gogithub.Repository {
	ctx := context.Background()
	repo, _, err := getClient().Repositories.Get(ctx, owner, name)
	if err != nil {
		panic(err)
	}
	// fmt.Println(repo)
	// fmt.Println("")
	// fmt.Println(response)

	return repo
}

/*

// Repository represents a GitHub repository.
type Repository struct {
	ID               *int             `json:"id,omitempty"`
	Owner            *User            `json:"owner,omitempty"`
	Name             *string          `json:"name,omitempty"`
	FullName         *string          `json:"full_name,omitempty"`
	Description      *string          `json:"description,omitempty"`
	Homepage         *string          `json:"homepage,omitempty"`
	DefaultBranch    *string          `json:"default_branch,omitempty"`
	MasterBranch     *string          `json:"master_branch,omitempty"`
	CreatedAt        *Timestamp       `json:"created_at,omitempty"`
	PushedAt         *Timestamp       `json:"pushed_at,omitempty"`
	UpdatedAt        *Timestamp       `json:"updated_at,omitempty"`
	HTMLURL          *string          `json:"html_url,omitempty"`
	Language         *string          `json:"language,omitempty"`
	Fork             *bool            `json:"fork"`
	ForksCount       *int             `json:"forks_count,omitempty"`
	NetworkCount     *int             `json:"network_count,omitempty"`
	OpenIssuesCount  *int             `json:"open_issues_count,omitempty"`
	StargazersCount  *int             `json:"stargazers_count,omitempty"`
	SubscribersCount *int             `json:"subscribers_count,omitempty"`
	WatchersCount    *int             `json:"watchers_count,omitempty"`
	Organization     *Organization    `json:"organization,omitempty"`
	Private           *bool   `json:"private"`
}

*/

func addRepoToDb(owner, name string) {
	repo := getBasicRepoInfo("flaviocopes", "gitometer-laravel")
	db.AddNewRepo(repo)

}

// /**
//  * Add the repo to the db
//  *
//  * @throws Exception if DB already contains the repository in `repositories` or `organizations_repositories`
//  * @return void
//  */
// public function addToDb()
// {
//     $data = $this->getBasicRepoInfo();

//     if ($this->repoExistsInDb($this->repository_owner, $this->repository_name)) {
//         Log::info('Repository already exists. Enable it.');

//         DB::table('repositories')
//             ->where('id', $this->id)
//             ->update([
//                 'enabled' => true
//             ]);
//     } else {
//         if ($data['fork'] === true) {
//             if (App::environment() !== 'testing') {
//                 //only allow forks in testing
//                 return;
//             }
//         }

//         $data['added_at'] = new \DateTime();
//         $data['enabled'] = true;

//         try {
//             $this->id = DB::table('repositories')->insertGetId($data);
//             Log::info('Added repository to repositories', $data);
//         } catch (Exception $e) {
//             Log::info('repositories already contains info', $data);
//         }
//     }

//     $organization_id = (new Organization($this->repository_owner))->getId();

//     try {
//         $organizations_repositories_data = [
//             'organization_id' => $organization_id,
//             'repository_id' => $this->id
//         ];

//         DB::table('organizations_repositories')->insert($organizations_repositories_data);
//         Log::info('Added repository to organizations_repositories', $organizations_repositories_data);
//     } catch (Exception $e) {
//         Log::info('organizations_repositories already contains info', $organizations_repositories_data);
//     }

//     Log::info('repository_added', [
//         'repository_owner' => $this->repository_owner,
//         'repository_name' => $this->repository_name,
//         'repository_id' => $this->id
//     ]);
// }
