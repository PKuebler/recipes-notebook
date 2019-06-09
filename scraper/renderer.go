package scraper

import (
	"encoding/json"
	"fmt"
)

type Renderer struct {
	Content Content
}

func NewRenderer(content Content) *Renderer {
	return &Renderer{
		Content: content,
	}
}

func (r *Renderer) RenderJson() string {
	output, _ := json.Marshal(r.Content)
	return string(output)
}

func (r *Renderer) RenderMarkdown() string {
	m := ""

	// Meta Infos
	m += "---\n"
	m += fmt.Sprintf("title: %s\n", r.Content.Heading)
	m += fmt.Sprintf("categories:\n")
	for _, cat := range r.Content.Categories {
		m += fmt.Sprintf("- %s\n", cat)
	}
	m += fmt.Sprintf("ingredients:\n")
	for _, ing := range r.Content.Ingredients {
		m += fmt.Sprintf("- %s\n", ing.Name)
	}
	m += "---\n\n"

	// Heading
	m += fmt.Sprintf("# %s\n\n", r.Content.Heading)

	// Ingredients
	m += fmt.Sprintf("## Zutaten\n\n")
	// max lengths
	leftCount := 0
	rightCount := 0

	for _, ing := range r.Content.Ingredients {
		if leftCount < len(ing.Amount) {
			leftCount = len(ing.Amount)
		}
		if rightCount < len(ing.Name) {
			rightCount = len(ing.Name)
		}
	}

	m += fmt.Sprintf("| %s | %s |\n", PadRight("Anzahl", " ", leftCount), PadRight("Zutaten", " ", rightCount))
	m += fmt.Sprintf("| %s | %s |\n", PadRight("", "-", leftCount), PadRight("", "-", rightCount))

	for _, ing := range r.Content.Ingredients {
		amount := PadLeft(ing.Amount, " ", leftCount)
		name := PadRight(ing.Name, " ", rightCount)

		m += fmt.Sprintf("| %s | %s |\n", amount, name)
	}
	m += "\n\n"

	// Instruction Text
	m += fmt.Sprintf("## Zubereitung\n\n")
	for _, step := range r.Content.InstructionText {
		m += fmt.Sprintf("%s\n\n", step)
	}

	m += fmt.Sprintf("Source: [%s](%s)", r.Content.URL, r.Content.URL)

	return m
}

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) > lenght {
			return str[len(str)-lenght : len(str)]
		}
	}
}
