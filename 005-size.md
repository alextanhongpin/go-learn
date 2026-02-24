```go
func formatSize(size int) string {
	var sizes = []string{"B", "KB", "MB", "GB", "TB"}

	f := float64(size)
	for i := len(sizes) - 1; i > -1; i-- {
		unit := math.Pow(1024, float64(i))
		if f > unit {
			return formatDecimal(f/unit) + " " + sizes[i]
		}
	}
	return formatDecimal(f) + " " + sizes[0]
}

func formatDecimal(num float64) string {
	// Format to 1 fixed decimal place.
	str := fmt.Sprintf("%.1f", num)

	// Trim trailing "0"s and then the trailing "." if present.
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")

	return str
}
```
