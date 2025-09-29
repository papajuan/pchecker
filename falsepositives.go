package pchecker

// DefaultFalsePositives is a list of words that may wrongly trigger the DefaultProfanities
var DefaultFalsePositives = map[string]bool{
	"analy":         true, // analysis, analytics
	"arsenal":       true,
	"assassin":      true,
	"assaying":      true, // was saying
	"assert":        true,
	"assign":        true,
	"assimil":       true,
	"assist":        true,
	"associat":      true,
	"assum":         true, // assuming, assumption, assumed
	"assur":         true, // assurance
	"banal":         true,
	"basement":      true,
	"bass":          true,
	"cass":          true, // cassie, cassandra, carcass
	"butter":        true, // butter, butterfly
	"button":        true,
	"canvass":       true,
	"circum":        true,
	"clitheroe":     true,
	"cockburn":      true,
	"cocktail":      true,
	"cumber":        true,
	"cumbing":       true,
	"cumulat":       true,
	"dickvandyke":   true,
	"document":      true,
	"evaluate":      true,
	"exclusive":     true,
	"expensive":     true,
	"explain":       true,
	"expression":    true,
	"grape":         true,
	"grass":         true,
	"glass":         true,
	"harass":        true,
	"hass":          true,
	"horniman":      true,
	"hotwater":      true,
	"identit":       true,
	"kassa":         true, // kassandra
	"kassi":         true, // kassie, kassidy
	"lass":          true, // class
	"leafage":       true,
	"libshitz":      true,
	"magnacumlaude": true,
	"mass":          true,
	"mocha":         true,
	"pass":          true, // compass, passion
	"penistone":     true,
	"phoebe":        true,
	"phoenix":       true,
	"pushit":        true,
	"sassy":         true,
	"saturday":      true,
	"scrap":         true, // scrap, scrape, scraping
	"serfage":       true,
	"sexist":        true, // systems exist, sexist
	"shoe":          true,
	"scunthorpe":    true,
	"shitake":       true,
	"stitch":        true,
	"sussex":        true,
	"therapist":     true,
	"therapeutic":   true,
	"tysongay":      true,
	"wass":          true,
	"wharfage":      true,
}
