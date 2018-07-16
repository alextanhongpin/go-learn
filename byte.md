For converting from a string to a byte slice, string -> []byte:
```go
[]byte(str)
```

For converting an array to a slice, [20]byte -> []byte:
```go
arr[:]
```

For copying a string to an array, string -> [20]byte:
```go
copy(arr[:], str)
```

Same as above, but explicitly converting the string to a slice first:
```go
copy(arr[:], []byte(str))
```
