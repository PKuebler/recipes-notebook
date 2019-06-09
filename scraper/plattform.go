package scraper

const (
	CHEFKOCH_URL = "chefkoch.de"
	LECKER_URL   = "lecker.de"
)

type Plattform string

const (
	PlattformChefkoch Plattform = CHEFKOCH_URL
	PlattformLecker   Plattform = LECKER_URL
)
