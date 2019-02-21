```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Connect to database.
	db, err := database.New(...)
	if err != nil {
		log.Fatal(err)
	}

	var b []byte
	stmt := `
		SELECT JSON_OBJECT(
			"id", HEX(id),
			"email", email,
			"email_verified", IF(email_verified = 1, true, false) IS true
		) FROM employee
	`
	err = db.QueryRow(stmt).Scan(&b)
	if err != nil {
		log.Fatal(err)
	}
	// Unmarshal only if you need to work with the data.
	// var m model.Employee
	// if err := json.Unmarshal(b, &m); err != nil {
	//         log.Fatal(err)
	// }
	// log.Println(m)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(b))
	})

	fmt.Println("listening to port *:4000")
	http.ListenAndServe(":4000", nil)
}
```

## With gin

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.New(...)
	if err != nil {
		log.Fatal(err)
	}

	var b []byte
	stmt := `
		SELECT JSON_OBJECT(
			"id", HEX(id),
			"email", email,
			"email_verified", IF(email_verified = 1, true, false) IS true
		) FROM employee
	`
	err = db.QueryRow(stmt).Scan(&b)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		var m map[string]interface{}
		json.Unmarshal(b, &m)
		// Send as string...
		// c.String(http.StatusOK, string(b))
		c.JSON(http.StatusOK, m)
	})
	r.Run()
}
```
