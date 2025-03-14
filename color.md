```go
// contrastColor calculates a contrasting color (either black or white) based on the input color's luminance.
func contrastColor(r, g, b int) string {
	// Counting the perceptive luminance - human eye favors green color...
	luminance := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255

	if luminance > 0.5 {
		return "black"
	}

	return "white"
}

// hexToRGBA converts a CSS hex color code to an RGBA color format.
func hexToRGBA(hex string) (int, int, int, float64, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 && len(hex) != 8 {
		return 0, 0, 0, 0, fmt.Errorf("invalid hex color format")
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 32)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	g, err := strconv.ParseInt(hex[2:4], 16, 32)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	b, err := strconv.ParseInt(hex[4:6], 16, 32)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var a int64 = 255
	if len(hex) == 8 {
		a, err = strconv.ParseInt(hex[6:8], 16, 32)
		if err != nil {
			return 0, 0, 0, 0, err
		}
	}

	return int(r), int(g), int(b), float64(a) / 255, nil
}
```
