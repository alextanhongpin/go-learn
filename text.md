# Stringcases
```go
package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var titleCaser = cases.Title(language.English, cases.NoLower)
var lowerCaser = cases.Lower(language.English)

var pascalTransformer = transform.Chain(
	norm.NFD,
	titleCaser,
	runes.Remove(&LetterOnly{}),
	norm.NFC,
)

var kebabTransformer = transform.Chain(
	norm.NFD,
	titleCaser,
	runes.Remove(&LetterOnly{}),
	&charBeforeUpper{SkipBeginning: true, Char: '-'},
	lowerCaser,
	norm.NFC,
)

var snakeTransformer = transform.Chain(
	norm.NFD,
	titleCaser,
	runes.Remove(&LetterOnly{}),
	&charBeforeUpper{SkipBeginning: true, Char: '_'},
	lowerCaser,
	norm.NFC,
)

func main() {
	input := []string{
		"userID",
		"hello world",
		"UserStore",
		"userAPI",
		"clientHTTP",
		// Doesn't work well.
		"user_store",
		"user_api",
		"JSONSerializer",
	}
	for _, t := range input {
		fmt.Println(t, toSnake(t), toCamel(t), toPascal(t), toKebab(t))
	}
	input = nil

	for _, s := range input {
		fmt.Println(ToPascal(s))
		fmt.Println(ToCamel(s))
		fmt.Println(ToSnake(s))
		fmt.Println(ToKebab(s))
		fmt.Println()
	}
}

func ToCamel(in string) string {
	out, _, err := transform.String(pascalTransformer, in)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(out[:1]) + out[1:]
}

func ToPascal(in string) string {
	out, _, err := transform.String(pascalTransformer, in)
	if err != nil {
		panic(err)
	}
	return out
}

func ToKebab(in string) string {
	out, _, err := transform.String(kebabTransformer, in)
	if err != nil {
		panic(err)
	}
	return out
}

func ToSnake(in string) string {
	out, _, err := transform.String(snakeTransformer, in)
	if err != nil {
		panic(err)
	}
	return out
}

type LetterOnly struct{}

func (l *LetterOnly) Contains(r rune) bool {
	return !unicode.IsLetter(r)
}

type charBeforeUpper struct {
	transform.NopResetter
	SkipBeginning bool
	Char          byte
}

func (t charBeforeUpper) skip(i int) bool {
	return !t.SkipBeginning || (t.SkipBeginning && i == 0)
}

func (t charBeforeUpper) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if !atEOF {
		return 0, 0, transform.ErrShortSrc
	}

	var j int
	k := len(src)

	for i, c := range src {
		left := src[i:]

		_, ok := commonInitialisms[string(left)]
		if ok {
			n := len(left)
			if k+1 >= len(dst) {
				return k + 1, i, transform.ErrShortSrc
			}
			dst[j] = t.Char
			k++
			j++
			for ii := i; ii < i+n; ii++ {
				dst[j] = src[ii]
				j++
			}
			break
		}
		if !t.skip(i) && c >= 'A' && c <= 'Z' {
			if k+1 >= len(dst) {
				return k + 1, i, transform.ErrShortSrc
			}
			dst[j] = t.Char
			dst[j+1] = c
			j++
			k++
		} else {
			dst[j] = c
		}
		j++
	}

	return k, len(src), nil
}

// https://github.com/golang/lint/blob/6edffad5e6160f5949cdefc81710b2706fbcd4f6/lint.go#LL766-L809
// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

var commonRe = regexp.MustCompile(`(ACL|API|ASCII|CPU|CSS|DNS|EOF|GUID|HTML|HTTP|HTTPS|ID|IP|JSON|LHS|QPS|RAM|RHS|RPC|SLA|SMTP|SQL|SSH|TCP|TLS|TTL|UDP|UI|UID|URI|URL|UTF8|UUID|VM|XML|XMPP|XSRF|XSS)`)
var repeatRe = regexp.MustCompile(`([^a-zA-Z]+)`)
var camelRe = regexp.MustCompile(`([a-z][A-Z])`)

func toSnake(s string) string {
	s = commonRe.ReplaceAllStringFunc(s, func(s string) string {
		// common initialisms can occur in the beginning as well as the end.
		return fmt.Sprintf("_%s_", strings.ToLower(s))
	})

	s = camelRe.ReplaceAllStringFunc(s, func(s string) string {
		return fmt.Sprintf("%s_%s", s[:1], strings.ToLower(s[1:]))
	})

	var sb strings.Builder
	for i, r := range s {
		if unicode.IsNumber(r) || unicode.IsLetter(r) {
			sb.WriteRune(unicode.ToLower(r))
		} else {
			if i != 0 && i != len(s)-1 {
				sb.WriteRune('_')
			}
		}
	}

	s = repeatRe.ReplaceAllStringFunc(sb.String(), func(s string) string {
		return s[:1]
	})

	return s
}

func toKebab(s string) string {
	s = toSnake(s)
	return strings.ReplaceAll(s, "_", "-")
}

func toPascal(s string) string {
	s = toSnake(s)

	words := strings.Split(s, "_")
	res := make([]string, len(words))
	for i, word := range words {
		wu := strings.ToUpper(word)
		if commonInitialisms[wu] {
			res[i] = wu
		} else {
			res[i] = titleCaser.String(word)
		}
	}

	return strings.Join(res, "")
}

func toCamel(s string) string {
	s = toSnake(s)

	words := strings.Split(s, "_")
	res := make([]string, len(words))
	for i, word := range words {
		if i == 0 {
			res[i] = word
			continue
		}
		wu := strings.ToUpper(word)
		if commonInitialisms[wu] {
			res[i] = wu
		} else {
			res[i] = titleCaser.String(word)
		}
	}

	return strings.Join(res, "")
}
```
