## Error message localization with golang validator

```go
package main

import (
	"fmt"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"

	"gopkg.in/go-playground/validator.v9"
	ja_translations "gopkg.in/go-playground/validator.v9/translations/ja"
)

type User struct {
	Username string `validate:"required"`
	Tagline  string `validate:"required,lt=10"`
}



func main() {
	ja := ja.New()
	uni := ut.New(ja, ja)
	trans, _ := uni.GetTranslator("ja")

	user := User{
		Username: "",
		Tagline:  "hello world, this is my playground",
	}
	validate := validator.New()
	ja_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(&user)
	if err != nil {
		// fmt.Println(err)
		// Translate all errors at once.
		errs := err.(validator.ValidationErrors)

		// Translate one. Produce a string.
		fmt.Println(errs[0].Translate(trans))

		// Translate all. Produce a map[string]string.
		// errsMap := errs.Translate(trans)
		// fmt.Println(errsMap["User.Tagline"])
	}
}
```

## References
https://phraseapp.com/blog/posts/internationalization-i18n-go/
