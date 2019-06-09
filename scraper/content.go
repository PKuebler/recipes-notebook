package scraper

type Ingredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type Content struct {
	URL             string       `json:"url"`
	Heading         string       `json:"heading"`
	Categories      []string     `json:"categories"`
	Ingredients     []Ingredient `json:"ingredients"`
	InstructionText []string     `json:"instruction_text"`
}

func NewContent(target string) Content {
	return Content{
		URL:             target,
		Categories:      []string{},
		Ingredients:     []Ingredient{},
		InstructionText: []string{},
	}
}
