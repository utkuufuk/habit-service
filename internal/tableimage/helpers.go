package tableimage

import (
	"fmt"
	"image/color"
	"strings"
)

func parseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255

	if len(x) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}

	if len(x) == 6 {
		fmt.Sscanf(x, "%02x%02x%02x", &r, &g, &b)
	}

	if len(x) == 8 {
		fmt.Sscanf(x, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}

	return r, g, b, a
}

func getColorByHex(hexColor string) color.RGBA {
	r, g, b, a := parseHexColor(hexColor)
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func wrapText(input string) (wrapped []string) {
	// split string into array of words
	words := strings.Fields(input)
	wordsLength := len(words)

	if wordsLength == 0 {
		return wrapped
	}

	var lineText string

	for i, word := range words {
		if len(lineText)+len(word)+1 >= wrapWordsLen {
			wrapped = append(wrapped, lineText)
			lineText = word
		} else {
			if lineText == "" {
				lineText += word
			} else {
				lineText += " " + word
			}

			// if it is the last word
			if i == wordsLength-1 {
				wrapped = append(wrapped, lineText)
			}

		}
	}

	return wrapped
}
