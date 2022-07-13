package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NERVEbing/douban-box-go/pkg/douban"
)

func main() {
	ghUsername := os.Getenv("GH_USER")
	ghToken := os.Getenv("GH_TOKEN")
	gistID := os.Getenv("GIST_ID")
	dbUser := os.Getenv("DOUBAN_USER")
	timezone := os.Getenv("TIMEZONE")

	gist, err := douban.NewGist(ghUsername, ghToken, dbUser, timezone)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = gist.BuildGitHubGist(ctx, gistID)
	if err != nil {
		panic(err)
	}

	fmt.Println("updating gist successfully")
}
