package douban

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mattn/go-runewidth"
)

const (
	MaxLines     = 5
	MaxLineWidth = 40
)

type gist struct {
	github    *github.Client
	rssDouBan *rssDouBan
}

func NewGist(ghUsername string, ghToken string, dbUser string) (*gist, error) {
	gist := &gist{}

	ghTransport := github.BasicAuthTransport{
		Username: strings.TrimSpace(ghUsername),
		Password: strings.TrimSpace(ghToken),
	}

	gist.github = github.NewClient(ghTransport.Client())
	rssDouBan, err := getRSSDouBan(dbUser)
	if err != nil {
		return nil, err
	}
	gist.rssDouBan = rssDouBan

	return gist, nil
}

func (g *gist) BuildGitHubGist(ctx context.Context, gistID string) error {
	ghGist, err := g.getGitHubGist(ctx, gistID)
	if err != nil {
		return err
	}

	f := g.getGitHubGistFile(ghGist)
	f.Content = github.String(g.formatRSSDouBanToGitHubGistContent())
	ghGist.Files[github.GistFilename(g.getGitHubGistFilename())] = f
	err = g.updateGitHubGist(ctx, gistID, ghGist)
	if err != nil {
		return err
	}

	return nil
}

func (g *gist) getGitHubGist(ctx context.Context, id string) (*github.Gist, error) {
	ghGist, _, err := g.github.Gists.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return ghGist, nil
}

func (g *gist) updateGitHubGist(ctx context.Context, id string, ghGist *github.Gist) error {
	_, _, err := g.github.Gists.Edit(ctx, id, ghGist)
	if err != nil {
		return err
	}

	return nil
}

func (g *gist) getGitHubGistFilename() string {
	return g.rssDouBan.Title
}

func (g *gist) getGitHubGistFile(ghGist *github.Gist) github.GistFile {
	return ghGist.Files[github.GistFilename(g.getGitHubGistFilename())]
}

func (g *gist) formatRSSDouBanToGitHubGistContent() string {
	content := ""
	spacesLen := 6

	max := g.getRSSDouBanItemsTitleMaxWidth()
	if max < MaxLineWidth {
		spacesLen = MaxLineWidth - max
	}

	for _, item := range g.rssDouBan.Items {
		spaces := strings.Repeat(" ", spacesLen)
		line := fmt.Sprintf("%s%s%s\n", item.PubDate, spaces, item.Title)
		content += line
	}

	return content
}

func (g *gist) getRSSDouBanItemsTitleMaxWidth() int {
	max := 0
	for _, item := range g.rssDouBan.Items {
		if max < runewidth.StringWidth(item.Title) {
			max = runewidth.StringWidth(item.Title)
		}
	}

	return max
}
