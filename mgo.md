# Working with mgo package in Go

## Connecting

```go
// DB holds the session to the mongo db
type DB struct {
	Session *mgo.Session
	Name    string
}

// New returns a new pointer to the DB struct
func New(addr, username, password, database, auth string) *DB {
	s, err := mgo.DialWithInfo(&mgo.DialInfo{
		Username: username,
		Password: password,
		Database: auth,
		Timeout:  time.Minute * 1,
		Addrs:    []string{addr},
	})
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)

	return &DB{
		Session: s,
		Name:    database,
	}
}
```

## Creating collection

```go
sess, err := mgo.DialWithInfo(&mgo.DialInfo{
  Username: username,
  Password: password,
  Database: auth,
  Timeout:  time.Minute * 1,
  Addrs:    []string{addr},
})

// Copy a new session
s := sess.Copy()

// Remember to close the session
defer s.Close()

// Do something with collection
c := s.DB("db_name").C("collection_name")

// or use the collection with the given session
c := sess.DB("db_name).C("collection_name").With(s)
```


## Setting index

To ensure if a field is unique:

```go
err := c.EnsureIndex(mgo.Index{
  Key:    []string{"login"},
  Unique: true,
})
```

## Single Upsert

```go
change, err := users.Upsert(
  bson.M{"login": "alextanhongpin"},
  bson.M{
    "$set": bson.M{
    // Remember to update the last updated date
      "updatedAt": time.Now().UTC().Format(time.RFC3339),
      "count":     10,
    },
    // Additionally you can include this field which will be inserted if
    // the document does not exist
    "$setOnInsert": bson.M {
      "createdAt": time.Now().UTC().Format(time.RFC3339),
    },
  },
  
)
```


## Bulk Upsert

Create a new document if it does not exist, and update an existing one with the given fields.

```go
// c is a *mgo.Collection
bulk := c.Bulk()

bulk.Upsert(bson.M{"login": "johndoe"}, bson.M{"$set": bson.M{"count": 1}})
bulk.Upsert(bson.M{"login": "alextanhongpin"}, bson.M{"$set": bson.M{"count": 10}})
bulk.Upsert(bson.M{"login": "hello"}, bson.M{"$set": bson.M{"count": 10}})

// Note that mgo can only support up to 1000 items for bulk. 
// Remember to partition the data you want to index.
change, err := bulk.Run()
```

## Sorting

Adding the `-` sign means sort by descending (big to small). 
```go
err = c.Find(nil).Sort("-timestamp").Limit(100).All(&res)
```

## Find One

```go
// Find  One
var user User
if err = c.Find(bson.M{"login": "alextanhongpin"}).
  Select(bson.M{"login": 0}). // Exclude login field
  One(&user); err != nil {
  log.Println(err)
}
log.Printf("user: %+v\n", user)
```

## Find All
```go
// Find All
var users []User
if err = c.Find(bson.M{}).All(&users); err != nil {
  log.Println(err)
}
```

## Remove All

```go
change, err := c.RemoveAll(bson.M{})
if err != nil {
  panic(err)
}
```
