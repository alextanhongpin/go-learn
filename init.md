
```go
// OK
// Pros: easy initialization
// Cons: only a single instance of client, not thread-safe
var instance *ethclient.Client

func GetClient(url string) *ethclient.Client {
	if client == nil {
		client, err := ethclient.Dial(url)
		if err != nil {
			log.Fatal(err)
		}
        instance = client
	}
	return instance
}
```

```go
// GOOD
// Reference: http://marcio.io/2015/07/singleton-pattern-in-go/
// Pros: easy initialization, thread-safe, initialize once
// Cons: only a single instance of client
var once sync.Once
var instance *ethclient.Client

func GetClient(url string)  *ethclient.Client  {
    once.Do(func() {
        client, err := ethclient.Dial(url)
        if err != nil {
            log.Fatal(err)
        }
        instance = client
    })
    return instance
}
```

```go
// BETTER
// Pros: Can invoke mutliple times with different params, thread-safe, no global variables
// Cons: Pass down through dependency injection, but it's a good practice

type Contract struct {
    Client  *ethclient.Client
}

func New(url string) (*ethclient.Client, err) {
    client, err := ethclient.Dial(url)
    if err != nil {
        return nil, err
    }
    return &Contract {
        Client: client,
    }, nil
}

// main.go
func main() {
    con, err := contract.New("http://localhost:8545")
    if err != nil {
        log.Fatal(err)
    }
    // can initialize multiple contracts with different address, not just one
    // con2, err := contract.New(addr2)
}```
