## Comparing strings effectively

tl;dr, use `strings.EqualFold` for comparison:

```go
// Good
if ok := strings.ToLower(a) == strings.ToLower(b); ok {}

// Better
if ok := strings.EqualFold(a, b); ok {}
```

References:
- https://www.digitalocean.com/community/questions/how-to-efficiently-compare-strings-in-go
