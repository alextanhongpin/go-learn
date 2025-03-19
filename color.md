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

# Using APCA

In the context of the Accessible Perceptual Contrast Algorithm (APCA), a minimum contrast value of Lc 75 is preferred for body text, while Lc 15 is considered the point of invisibility for many users, especially for thin lines or borders.

```go
package video_generator

import "math"

func sRGBtoY(srgb [3]int) float64 {
	r := math.Pow(float64(srgb[0])/255, 2.4)
	g := math.Pow(float64(srgb[1])/255, 2.4)
	b := math.Pow(float64(srgb[2])/255, 2.4)
	y := 0.2126729*r + 0.7151522*g + 0.0721750*b

	if y < 0.022 {
		y += math.Pow(0.022-y, 1.414)
	}
	return y
}

func contrast(fg, bg [3]int) float64 {
	yfg := sRGBtoY(fg)
	ybg := sRGBtoY(bg)
	c := 1.14

	if ybg > yfg {
		c *= math.Pow(ybg, 0.56) - math.Pow(yfg, 0.57)
	} else {
		c *= math.Pow(ybg, 0.65) - math.Pow(yfg, 0.62)
	}

	if math.Abs(c) < 0.1 {
		return 0
	} else if c > 0 {
		c -= 0.027
	} else {
		c += 0.027
	}

	return c * 100
}

func colorContrast(foreground, background [3]int) string {
	// The value can be negative if you have a white foreground over black background
	if math.Abs(contrast(foreground, background)) >= 75 {
		return "white"
	}

	return "black"
}
```
