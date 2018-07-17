## Initializing vgo

This will create a `go.mod` file:

```bash
$ vgo mod -init
```

## Vendor

```bash
$ vgo mod -vendor

# To view help
$ vgo help mod
```

## Pulling specific version
```
# Pull v1.0.0 from Github

$ vgo get github.com/alextanhongpin/go-vgo-test
vgo: finding github.com/alextanhongpin/go-vgo-test v1.0.0
vgo: downloading github.com/alextanhongpin/go-vgo-test v1.0.0

$ vgo run main.go
hello, v1!
```

```
# Pull v2.0.0 from Github

$ vgo get github.com/alextanhongpin/go-vgo-test@v2.0.0
vgo: finding github.com/alextanhongpin/go-vgo-test v2.0.0
vgo: downloading github.com/alextanhongpin/go-vgo-test v0.0.0-20180717013519-63d87eea745d

$ vgo run main.go
vgo: finding github.com/alextanhongpin/go-vgo-test v0.0.0-20180717013519-63d87eea745d
hello, v2!```
