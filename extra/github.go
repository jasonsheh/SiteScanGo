package extra

import (
	"../utils"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
)

func Search(keyword string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "76ce8f923e2448ee6c2f288138bd36d02f4600c5"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	codeResults, _, err := client.Search.Code(ctx, keyword, nil)
	utils.CheckError(err)

	fmt.Println(*codeResults.Total)
	codeResult := codeResults.CodeResults
	for _, result := range codeResult {
		url := strings.Replace(*result.HTMLURL, "/blob", "", 1)
		url = strings.Replace(url, "github.com", "raw.githubusercontent.com", 1)
		fmt.Println(url)
	}
}

func GithubLeak(domainKeyword string){
	keywordList := []string{"password="}

	for _, keyword := range keywordList {
		keyword = keyword + `+"`+domainKeyword+`"`
		Search(keyword)
	}
}
