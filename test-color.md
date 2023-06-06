# Adding color output for tests

## Easiest

Install this package: https://github.com/rakyll/gotest

Then replace `go test` with `gotest`.

<img width="789" alt="image" src="https://github.com/alextanhongpin/go-learn/assets/6033638/fe2917a3-a871-429f-92b5-4dc72c0023c7">


## Harder
You can use grc, a generic colourizer, to colourize anything.

On Debian/Ubuntu, install with `apt-get install grc`. On a Mac with , `brew install grc`.

Create a config directory in your home directory:

```bash
mkdir ~/.grc
```

Then create your personal grc config in `~/.grc/grc.conf`:

```
# Go
\bgo.* test\b
conf.gotest
```

Then create a Go test colourization config in `~/.grc/conf.gotest`, such as:

```
regexp==== RUN .*
colour=blue
-
regexp=--- PASS: .*
colour=green
-
regexp=^PASS$
colour=green
-
regexp=^(ok|\?) .*
colour=magenta
-
regexp=--- FAIL: .*
colour=red
-
regexp=[^\s]+\.go(:\d+)?
colour=cyan
-
regexp=\s\-\s.*
colour=red
-
regexp=Snapshot\(\-\)
colour=red
-
regexp=\s\+\s.*
colour=green
-
regexp=Received\(\+\)
colour=green
-
```

I customize it to show the diff `+/-`, especially useful with library like `go-cmp`.

Now you can run Go tests with:

```
grc go test -v ./..
```

Sample output:

<img width="898" alt="image" src="https://github.com/alextanhongpin/go-learn/assets/6033638/36578e82-fec6-4534-9c02-734b5e356c24">


To avoid typing grc all the time, add an alias to your shell (if using Bash, either `~/.bashrc` or `~/.bash_profile` or both, depending on your OS):

```
alias go=grc go
```

Now you get colourization simply by running:

```
go test -v ./..
```

Reference here: https://stackoverflow.com/questions/27242652/colorizing-golang-test-run-output
