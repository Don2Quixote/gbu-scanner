package blog

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gbu-scanner/internal/entity"

	"gbu-scanner/pkg/logger"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// Blog is implementation for scanner.Blog interface
type Blog struct {
	host       string // Host where blog hosted (go.dev)
	blogPath   string // Path to all posts (/blog/all)
	protocol   string // Protocol to use (http or https)
	httpClient HTTPClient
	log        logger.Logger
}

// New returns scanner.Blog implementation
func New(host string, blogPath string, https bool, client HTTPClient, log logger.Logger) *Blog {
	protocol := "http"
	if https {
		protocol = "https"
	}
	return &Blog{
		host:       host,
		blogPath:   blogPath,
		protocol:   protocol,
		httpClient: client,
		log:        log,
	}
}

func (p *Blog) GetPosts(ctx context.Context) ([]entity.Post, error) {
	url := fmt.Sprintf("%s://%s%s", p.protocol, p.host, p.blogPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create request")
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't do request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("response status code is not OK (%s)", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't get document from response")
	}

	elements := doc.Find(".blogtitle")
	posts := make([]entity.Post, 0, elements.Length())
	elements.Each(func(i int, blogtitle *goquery.Selection) {
		linkElement := blogtitle.Find("a")
		if linkElement.Length() == 0 {
			p.log.Error("no link element")
			return
		}

		title := strings.TrimSpace(linkElement.Text())

		dateElement := linkElement.Next()
		if dateElement.Length() == 0 {
			p.log.Error("no next element after link")
			return
		}

		date, err := time.Parse(dateLayout, strings.TrimSpace(dateElement.Text()))
		if err != nil {
			p.log.Error("can't parse date")
			return
		}

		authorElement := blogtitle.Find(".author")
		if authorElement.Length() == 0 {
			p.log.Error("no author element")
			return
		}

		author := strings.TrimSpace(authorElement.Text())

		summaryElement := blogtitle.Next()
		if summaryElement.Length() == 0 {
			p.log.Error("no next element after blogtitle")
			return
		}

		summary := strings.TrimSpace(summaryElement.Text())

		path, ok := linkElement.Attr("href")
		if !ok {
			p.log.Error("no href attribute on link element")
			return
		}
		url := fmt.Sprintf("%s://%s%s", p.protocol, p.host, path)

		posts = append(posts, entity.Post{
			Title:   title,
			Date:    date,
			Author:  author,
			Summary: summary,
			URL:     url,
		})
	})

	return posts, nil
}
