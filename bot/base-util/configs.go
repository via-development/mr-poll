package baseUtil

import (
	"os"
	"strconv"
)

type config struct {
	EmbedColor int
}

var configs = map[int]config{
	0: { // Production Config
		EmbedColor: 0xFF5C40,
	},
	1: { // Developer Config
		EmbedColor: 0xFF5C40,
	},
}

var Config config

func init() {
	if len(os.Args) < 2 {
		panic("You must input a config profile ID.")
	}

	ConfigId, err := strconv.Atoi(os.Args[1])

	if err != nil {
		panic(err)
	}

	Config = configs[ConfigId]
}
