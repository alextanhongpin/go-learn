https://golang.org/doc/effective_go#Getters

- the problem with private fields
- it makes constructor lengthy when you habe moe fields
- it becomes hard to set value in another package
- to avoid lengthy constructor, a builder is required, but that means more code
- best to have all public name, setter only when the field has business logic
