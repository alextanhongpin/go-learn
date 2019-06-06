Better golang
- For function that needs to be initialised , we can prefix `Must` e.g. database.MustConnect, and just `log.Fatal` the error
- Create `database.Env()` to get the values from the environment variable. Since there are a lot of different ways of getting config, we use `Env` to indicate that they are required in the environment variables
- No *pointer slice please, there are reason when to use pointer slice (modifying the slice in place when it is really large), but for most cases, this can be avoided.
