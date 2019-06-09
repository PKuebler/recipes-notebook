# Recipes Notebook

Store recipes from different platforms in your own notebook on a Github or Gitlab. It should be usable as AWS Lambda function.

The JSON or Markdown files can subsequently be used e.g. Rendered by Static Site Generator or exported as a PDF.

## Example

```go
package main

import (
	"fmt"

	"github.com/gosimple/slug"

	"github.com/PKuebler/recipes-notebook/publish"
	"github.com/PKuebler/recipes-notebook/scraper"
)

func main() {
	token := "XXXXXXXXXXXXXXXXXX" // Gitlab Token

	//	s := NewScraper("https://www.chefkoch.de/rezepte/1770241286984131/Mediterrane-Bohnen-Nudel-Pfanne.html")
	s := scraper.NewScraper("https://www.lecker.de/20-minuten-maissuppe-74200.html")
	s.Crawl()

	r := scraper.NewRenderer(s.Content)

	fmt.Println(r.RenderMarkdown())
	fmt.Println(r.RenderJson())

	name := slug.MakeLang(s.Content.Heading, "de")

	git := publish.NewGitLab("000000000") // Project ID

	// Add Files to Commit
	git.AddFile(fmt.Sprintf("content/recipes/%s.md", name), r.RenderMarkdown())
	git.AddFile(fmt.Sprintf("static/json/%s.json", name), r.RenderJson())

	// Login
	git.Auth(token)

	// Commit
	fmt.Println(git.Commit("master", fmt.Sprintf("Add %s", s.Content.Heading)))
}

```

## Inspiration

- [Chef](https://github.com/runepiper/chef)(PHP) by [Runepiper](https://github.com/runepiper)

## ToDo

- Add Github API
- Extend code to simplify add new platforms