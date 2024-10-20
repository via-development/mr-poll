package main

import (
	"os"
	"strings"
)

func main() {
	files, err := os.ReadDir("./prisma")
	if err != nil {
		panic(err)
	}

	combinedFileContent := ""

	for _, file := range files {
		filename := file.Name()
		if file.IsDir() || filename == "schema.prisma" || !strings.HasSuffix(filename, ".prisma") {
			continue
		}

		var content string
		{
			c, err := os.ReadFile("./prisma/" + filename)
			if err != nil {
				panic(err)
			}
			content = string(c)
		}

		content = strings.TrimSpace(content)
		if filename == "base.prisma" {
			combinedFileContent = content + combinedFileContent
		} else {
			combinedFileContent += "\n\n" + content
		}

	}

	err = os.WriteFile("./prisma/schema.prisma", []byte(combinedFileContent), 0666)
	if err != nil {
		panic(err)
	}
}
