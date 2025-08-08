package util

import (
	"os"
	"strconv"
)

type config struct {
	EmbedColor int
	ShardCount int
	ShardIds   []int
}

var configs = map[int]config{
	0: { // Developer Config
		EmbedColor: 0x40FFAC,
		ShardCount: 1,
		ShardIds:   []int{0},
	},
	1: { // Production Config
		EmbedColor: 0x40FFAC,
		ShardCount: 1,
		ShardIds:   []int{0, 1, 2},
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
