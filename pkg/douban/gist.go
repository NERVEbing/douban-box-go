package douban

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v58/github"
	"github.com/mattn/go-runewidth"
)

const (
	MaxLineWidth = 40
)

type Gist struct {
	github    *github.Client
	rssDouBan *rssDouBan
}

func NewGist(ghUsername string, ghToken string, dbUser string, timezone string) (*Gist, error) {
	gist := &Gist{}

	ghTransport := github.BasicAuthTransport{
		Username: strings.TrimSpace(ghUsername),
		Password: strings.TrimSpace(ghToken),
	}

	gist.github = github.NewClient(ghTransport.Client())
	rssDouBan, err := getRSSDouBan(dbUser, timezone)
	if err != nil {
		return nil, err
	}
	gist.rssDouBan = rssDouBan

	return gist, nil
}

func (g *Gist) BuildGitHubGist(ctx context.Context, gistID string) error {
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

func (g *Gist) getGitHubGist(ctx context.Context, id string) (*github.Gist, error) {
	ghGist, _, err := g.github.Gists.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return ghGist, nil
}

func (g *Gist) updateGitHubGist(ctx context.Context, id string, ghGist *github.Gist) error {
	_, _, err := g.github.Gists.Edit(ctx, id, ghGist)
	if err != nil {
		return err
	}

	return nil
}

func (g *Gist) getGitHubGistFilename() string {
	return g.rssDouBan.Title
}

func (g *Gist) getGitHubGistFile(ghGist *github.Gist) github.GistFile {
	return ghGist.Files[github.GistFilename(g.getGitHubGistFilename())]
}

func (g *Gist) formatRSSDouBanToGitHubGistContent() string {
	content := ""

	for _, item := range g.rssDouBan.Items {
		lineWidth := runewidth.StringWidth(item.Stars + item.Title)
		if lineWidth > MaxLineWidth {
			continue
		}
		line := fmt.Sprintf("%s   %s   %s\n", item.PubDate, item.Stars, item.Title)
		content += line
	}

	return content
}
