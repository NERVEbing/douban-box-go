package douban

import (
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type rssDouBan struct {
	Title string
	Items []*rssDouBanItem
}

type rssDouBanItem struct {
	Title   string
	PubDate string
	Stars   string
}

const (
	CSTLayout = "2006.01.02"
)

func getRSSDouBan(dbUserID string) (*rssDouBan, error) {
	gist := &rssDouBan{}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.douban.com/feed/people/" + dbUserID + "/interests")
	if err != nil {
		return nil, err
	}

	gist.Title = formatRSSTitle(feed.Title)

	for _, item := range feed.Items {
		gistItem := &rssDouBanItem{
			Title:   formatRSSItemTitle(item.Title),
			PubDate: formatRSSItemDate(item.PublishedParsed),
			Stars:   formatRSSItemStars(item.Description),
		}
		gist.Items = append(gist.Items, gistItem)
	}

	return gist, nil
}

func formatRSSTitle(s string) string {
	if strings.Contains(s, " ") {
		arr := strings.Split(s, " ")
		if len(arr) != 2 {
			return s
		}
		userName := arr[0]
		s = fmt.Sprintf("Douban@%s", userName)
	}

	return s
}

func formatRSSItemTitle(s string) string {
	r := []rune(s)
	s1 := string(r[0:2])
	s2 := string(r[2:])

	switch {
	case strings.Contains(s1, "读"):
		s1 = "📖 " + s1
	case strings.Contains(s1, "看"):
		s1 = "🎬 " + s1
	case strings.Contains(s1, "听"):
		s1 = "🎵 " + s1
	default:
		s1 = "🧐" + s1
	}

	s = fmt.Sprintf("%s《%s》", s1, s2)

	return s
}

func formatRSSItemDate(t *time.Time) string {
	return t.Add(time.Hour * 8).Format(CSTLayout)
}

func formatRSSItemStars(s string) string {
	str := ""
	tag := "推荐: "
	y := "♥"
	n := "♡"
	x := "⚇"

	switch {
	case strings.Contains(s, tag+"很差"):
		str = strings.Repeat(y, 1) + strings.Repeat(n, 4)
	case strings.Contains(s, tag+"较差"):
		str = strings.Repeat(y, 2) + strings.Repeat(n, 3)
	case strings.Contains(s, tag+"还行"):
		str = strings.Repeat(y, 3) + strings.Repeat(n, 2)
	case strings.Contains(s, tag+"推荐"):
		str = strings.Repeat(y, 4) + strings.Repeat(n, 1)
	case strings.Contains(s, tag+"力荐"):
		str = strings.Repeat(y, 5) + strings.Repeat(n, 0)
	default:
		str = strings.Repeat(x, 5)
	}

	return str
}
