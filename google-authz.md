## Inspecting Google IdToken on the server side

Using the `oauth2` package:

```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/oauth2/v2"
)

func main() {
	idToken := "the idToken you want to inspect"
	
	oauth2Service, err := oauth2.New(&http.Client{})
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tokenInfoCall.Context(ctx)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	spew.Dump(tokenInfo)
}

/*
(*oauth2.Tokeninfo)(0xc000ff4140)({
  Audience: (string) (len=72) "<google-client-id>.apps.googleusercontent.com",
  Email: (string) (len=23) "john.doe@gmail.com",
  ExpiresIn: (int64) 3599,
  IssuedTo: (string) (len=72) "<google-client-id>.apps.googleusercontent.com",
  Scope: (string) "",
  UserId: (string) (len=21) "<user-id>",
  VerifiedEmail: (bool) true,
  ServerResponse: (googleapi.ServerResponse) {
    HTTPStatusCode: (int) 200,
    Header: (http.Header) (len=11) {}
    // Not shown...
  }
})
*/
```


Using the [google-api-go-client](https://github.com/googleapis/google-api-go-client) package:

```go
package main

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/idtoken"
)

func main() {
	var (
		googleClientID = "<google-client-id>"
		token          = "the idToken you want to inspect"
	)

	ctx := context.WihtTimeout(context.Background(), 10*time.Second)
	payload, err := idtoken.Validate(ctx, token, googleClientID)
	if err != nil {
		return nil, err
	}

	spew.Dump(payload)
}

/*
(*idtoken.Payload)(0xc000fe6370)({
  Issuer: (string) (len=19) "accounts.google.com",
  Audience: (string) (len=72) "<google-client-id>.apps.googleusercontent.com",
  Expires: (int64) 1590898001,
  IssuedAt: (int64) 1590894401,
  Subject: (string) (len=21) "<user-id>",
  Claims: (map[string]interface {}) <nil>
})
*/
```

## Extracting information from id token

NOTE: You need to verify the id token first with the method above.

```go
package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
)

type TokenInfo struct {
	Iss string `json:"iss"`
	// userId
	Sub string `json:"sub"`
	Azp string `json:"azp"`
	// clientId
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
	// expired time
	Exp int64 `json:"exp"`

	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Local         string `json:"locale"`
	jwt.StandardClaims
}

func main() {
	var (
		token = "the idToken you want to inspect"
	)

	jwtToken, _, err := new(jwt.Parser).ParseUnverified(token, &TokenInfo{})
	if err != nil {
		log.Fatal(err)
	}
	tokenInfo, ok := jwtToken.Claims.(*TokenInfo)
	if !ok {
		log.Fatal("invalid idToken")
	}
	spew.Dump(tokenInfo)

}

/*
(*google.TokenInfo)(0xc000ff2c60)({
  Iss: (string) (len=19) "accounts.google.com",
  Sub: (string) (len=21) "<user-id>",
  Azp: (string) (len=72) "<google-client-id>.apps.googleusercontent.com",
  Aud: (string) (len=72) "<google-client-id>.apps.googleusercontent.com",
  Iat: (int64) 1590894401,
  Exp: (int64) 1590898001,
  Email: (string) (len=23) "john.doe@gmail.com",
  EmailVerified: (bool) true,
  AtHash: (string) (len=22) "xyz",
  Name: (string) (len=8) "John Doe",
  GivenName: (string) (len=4) "John",
  FamilyName: (string) (len=3) "Doe",
  Picture: (string) (len=89) "https://imgurl",
  Local: (string) (len=2) "en",
  StandardClaims: (jwt.StandardClaims) {
    Audience: (string) "",
    ExpiresAt: (int64) 0,
    Id: (string) (len=40) "<user-id>",
    IssuedAt: (int64) 0,
    Issuer: (string) "",
    NotBefore: (int64) 0,
    Subject: (string) ""
  }
})

*/
```
