What if we have a function with side effects and want to create a reusable one? 

We can do sth like sql.Open.

We register a predefined struct that fulfils the interface, and open the given name to get the impelem station. 


So for a repository, we can open an in memory or postgres. 

```
repo.Register(name, &Driver{})
```
