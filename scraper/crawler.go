package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type Scraper struct {
	Plattform Plattform
	URL       string
	Content   Content
}

func NewScraper(target string) *Scraper {
	u, _ := url.Parse(target)
	host_parts := strings.Split(u.Host, ".")

	domain := u.Host

	if len(host_parts) > 2 {
		domain = strings.Join(host_parts[len(host_parts)-2:], ".")
	}

	fmt.Println(domain)
	plattform := PlattformChefkoch

	switch domain {
	case CHEFKOCH_URL:
		plattform = PlattformChefkoch
	case LECKER_URL:
		plattform = PlattformLecker
	}

	return &Scraper{
		URL:       target,
		Content:   NewContent(target),
		Plattform: plattform,
	}
}

func (s *Scraper) Crawl() {
	c := colly.NewCollector(
		colly.UserAgent("myUserAgent"),
		colly.MaxDepth(1),
		colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		s.CrawlHeading(e)
		s.CrawlCategories(e)
		s.CrawlInstructionText(e)
		s.CrawlIncredients(e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finish")
		fmt.Println(s.Content.Heading)
		fmt.Println(s.Content.Ingredients)
		fmt.Println(s.Content.Categories)
		fmt.Println(s.Content.InstructionText)
	})

	c.Visit(s.URL)
}

func (s *Scraper) CrawlHeading(e *colly.HTMLElement) {
	switch s.Plattform {
	case PlattformChefkoch:
		s.Content.Heading = e.ChildText(".page-title")
	case PlattformLecker:
		s.Content.Heading = e.ChildText(".article.recipe .article-header h1")
	}

	s.Content.Heading = s.clean(s.Content.Heading)
}

func (s *Scraper) CrawlCategories(e *colly.HTMLElement) {
	query := ".tagcloud li"

	switch s.Plattform {
	case PlattformLecker:
		query = ".list.list--tags li"
	}

	// Categories
	e.ForEach(query, func(_ int, el *colly.HTMLElement) {
		cat := s.clean(el.Text)
		s.Content.Categories = append(s.Content.Categories, cat)
	})
}

func (s *Scraper) CrawlInstructionText(e *colly.HTMLElement) {
	switch s.Plattform {
	case PlattformChefkoch:
		text := s.clean(e.ChildText(".instructions"))
		s.Content.InstructionText = []string{text}
	case PlattformLecker:
		e.ForEach(".list.list--preparation dd", func(_ int, el *colly.HTMLElement) {
			text := s.clean(el.Text)
			s.Content.InstructionText = append(s.Content.InstructionText, text)
		})
	}
}

func (s *Scraper) CrawlIncredients(e *colly.HTMLElement) {
	query := ".incredients tbody tr"
	queryAmount := "td:first-child"
	queryName := "td:nth-child(2)"

	switch s.Plattform {
	case PlattformLecker:
		query = ".list.list--ingredients li"
		queryAmount = ".quantityBlock"
		queryName = ".ingredientBlock"
	}

	e.ForEach(query, func(_ int, el *colly.HTMLElement) {
		amount := s.clean(el.ChildText(queryAmount))
		name := s.clean(el.ChildText(queryName))

		s.Content.Ingredients = append(s.Content.Ingredients, Ingredient{
			Name:   name,
			Amount: amount,
		})
	})
}

func (s *Scraper) clean(text string) string {
	insideRegEx := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	text = strings.TrimSpace(text)
	return insideRegEx.ReplaceAllString(text, " ")
}
