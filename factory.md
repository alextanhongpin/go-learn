# Attempt to make generic nested factory
```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

var mappings = map[string]map[string]any{}
var mu sync.Mutex

func Init(ctx context.Context) context.Context {
	// TODO: Cleanup with t.Cleanup
	return withGeneratedID(ctx)
}
func loadAll(ctx context.Context) map[string]any {
	id, ok := generatedID(ctx)
	if !ok {
		panic("invalid id")
	}
	mu.Lock()
	defer mu.Unlock()
	return mappings[id]
}

func loadOne[T any](ctx context.Context, prefix, id string) T {
	id, ok := generatedID(ctx)
	if !ok {
		panic("invalid id")
	}
	mu.Lock()
	defer mu.Unlock()
	return mappings[id][fmt.Sprintf("%s:%s", prefix, id)].(T)
}

func store(ctx context.Context, prefix, id string, val any) {
	id, ok := generatedID(ctx)
	if !ok {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	_, ok = mappings[id]
	if !ok {
		mappings[id] = make(map[string]any)
	}
	mappings[id][fmt.Sprintf("%s:%s", prefix, id)] = val
}

func main() {
	p := &Product{}
	t := reflect.ValueOf(p).Elem()
	t.FieldByName("ID").Set(reflect.ValueOf(1))
	fmt.Println(t.FieldByName("Haha").IsValid())
	fmt.Println(p)

	ctx := context.Background()
	actx := Init(ctx)
	p = createProduct(actx, map[string]any{
		"Product.Expensive": true,
	})
	fmt.Println(p)
	fmt.Println(loadOne[*Product](actx, "Product", fmt.Sprint(p.ID)))
	fmt.Println(loadAll(actx))

	fmt.Println(createProductCategory(ctx))
	fmt.Println(createProductCategorySubcategory(ctx))

	bctx := Init(ctx)
	p = createProduct(bctx, map[string]any{
		// entity.field to change the value
		// entity: variant to decide on the variant
		"Product.Name":              "toy car",
		"ProductCategory.Name":      "toys",
		"Product.WithSubcategories": 10,
	})
	fmt.Println(p)
	fmt.Println(loadAll(bctx))
	fmt.Println(loadOne[[]ProductCategorySubcategory](bctx, "ProductCategorySubcategoryListByProductID", fmt.Sprint(p.ID)))
}

// What is a better kv naming? we need to prefix it, e.g. Product.ID since we are passing it down multiple layers.
type KV map[string]any

type Product struct {
	ID    int
	Name  string
	Price int
}

type ProductCategory struct {
	ID        int
	Name      string
	ProductID int
}

type ProductCategorySubcategory struct {
	ID                int
	ProductCategoryID int
	Name              string
}

func buildProduct(kvs ...KV) *Product {
	p := Product{
		ID:    0,
		Name:  "chair",
		Price: 10,
	}

	for _, kv := range kvs {
		for k, v := range kv {
			if !strings.HasPrefix(k, "Product.") {
				continue
			}
			switch k {
			case "Product.Expensive":
				p.Price = 999_999
			// Avoid setting fields like this individually, rather, start with a scenario.
			case "Product.Price":
				p.Price = v.(int)
			case "Product.Name":
				p.Name = v.(string)
			case "Product.WithSubcategories":
			default:
				panic("invalid product variant: " + k)
			}
		}
	}

	return &p
}

func createProduct(ctx context.Context, kvs ...KV) *Product {
	p := buildProduct(kvs...)
	p.ID = 123
	fmt.Println("creating product", p)
	// AFTER
	for _, kv := range kvs {
		for k, v := range kv {
			if !strings.HasPrefix(k, "Product.") {
				continue
			}
			switch k {
			case "Product.WithSubcategories":
				pc := createProductCategory(ctx, append(kvs, map[string]any{
					"ProductCategory.ProductID": p.ID,
				})...)
				n := v.(int)
				subcategories := make([]ProductCategorySubcategory, n)
				for i := 0; i < n; i++ {
					subcategories[i] = *createProductCategorySubcategory(ctx, append(kvs, map[string]any{
						"ProductCategorySubcategory.ProductCategoryID": pc.ID,
						"ProductCategorySubcategory.Name":              fmt.Sprintf("subcategory_%d", i),
					})...)
				}
				store(ctx, "ProductCategorySubcategoryListByProductID", fmt.Sprint(p.ID), subcategories)
			}
		}
	}
	store(ctx, "Product", fmt.Sprint(p.ID), p)
	return p
}

func buildProductCategory(kvs ...KV) *ProductCategory {
	p := ProductCategory{
		ID:        0,
		Name:      "furniture",
		ProductID: 0,
	}

	for _, kv := range kvs {
		for k, v := range kv {
			if !strings.HasPrefix(k, "ProductCategory.") {
				continue
			}
			switch k {
			case "ProductCategory.ProductID":
				p.ProductID = v.(int)
			case "ProductCategory.Name":
				p.Name = v.(string)
			default:
				panic("invalid product category variant: " + k)
			}
		}
	}

	return &p
}

func createProductCategory(ctx context.Context, kvs ...KV) *ProductCategory {
	p := buildProductCategory(kvs...)
	if p.ProductID == 0 {
		p.ProductID = createProduct(ctx, kvs...).ID
	}
	p.ID = 65
	fmt.Println("creating product category", p)
	store(ctx, "ProductCategory", fmt.Sprint(p.ID), p)
	return p
}

func buildProductCategorySubcategory(kvs ...KV) *ProductCategorySubcategory {
	p := &ProductCategorySubcategory{
		ID:                0,
		Name:              "bed and tables",
		ProductCategoryID: 0,
	}

	for _, kv := range kvs {
		for k, v := range kv {
			typ, field, _ := strings.Cut(k, ".")
			if typ != "ProductCategorySubcategory" {
				continue
			}
			e := reflect.ValueOf(p).Elem()
			f := e.FieldByName(field)
			if !f.IsValid() {
				fmt.Println("invalid field name", field)
				continue
			}
			f.Set(reflect.ValueOf(v))
		}
	}

	return p
}

func createProductCategorySubcategory(ctx context.Context, kvs ...KV) *ProductCategorySubcategory {
	p := buildProductCategorySubcategory(kvs...)
	if p.ProductCategoryID == 0 {
		p.ProductCategoryID = createProductCategory(ctx, kvs...).ID
	}
	p.ID = 23
	fmt.Println("creating product category subcategory", p)
	store(ctx, "ProductCategorySubcategory", fmt.Sprint(p.ID), p)
	return p
}

// Don't do this, create it under createProduct instead.
func createProductCategorySubcategoryList() {}

type contextKey string

var factoryContextKey contextKey = "factory"

func withGeneratedID(ctx context.Context) context.Context {
	_, ok := generatedID(ctx)
	if ok {
		return ctx
	}
	return context.WithValue(ctx, factoryContextKey, fmt.Sprint(time.Now()))
}

func generatedID(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(factoryContextKey).(string)
	return s, ok
}
```


## Using reflect

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
	"log"
	"maps"
	"reflect"
	"sync"
	"time"
)

var mappings = map[string]map[string]any{}
var mu sync.Mutex

func Init(ctx context.Context) context.Context {
	// TODO: Cleanup with t.Cleanup
	return withGeneratedID(ctx)
}
func loadAll(ctx context.Context) map[string]any {
	id, ok := generatedID(ctx)
	if !ok {
		panic("invalid id")
	}
	mu.Lock()
	defer mu.Unlock()
	return mappings[id]
}

func loadOne[T any](ctx context.Context, prefix, id string) T {
	id, ok := generatedID(ctx)
	if !ok {
		panic("invalid id")
	}
	mu.Lock()
	defer mu.Unlock()
	return mappings[id][fmt.Sprintf("%s:%s", prefix, id)].(T)
}

func store(ctx context.Context, prefix, id string, val any) {
	id, ok := generatedID(ctx)
	if !ok {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	_, ok = mappings[id]
	if !ok {
		mappings[id] = make(map[string]any)
	}
	mappings[id][fmt.Sprintf("%s:%s", prefix, id)] = val
}

var productExpensive = KV{
	"Price": 999_999,
}

func main() {
	p := &Product{}
	t := reflect.ValueOf(p).Elem()
	t.FieldByName("ID").Set(reflect.ValueOf(1))
	fmt.Println(t.FieldByName("Haha").IsValid())
	fmt.Println(p)

	ctx := context.Background()
	actx := Init(ctx)
	p = createProduct(actx, productExpensive)
	fmt.Println(p)
	fmt.Println(loadOne[*Product](actx, "Product", fmt.Sprint(p.ID)))
	fmt.Println(loadAll(actx))

	fmt.Println(createProductCategory(ctx))
	fmt.Println(createProductCategorySubcategory(ctx))

	bctx := Init(ctx)
	p = createProduct(bctx, KV{
		// entity.field to change the value
		// entity: variant to decide on the variant
		"Name": "toy car",
		"Category": KV{
			"Name": "toys",
		},
		"WithSubcategories": 10,
	})
	fmt.Println(p)
	fmt.Println(loadAll(bctx))
	fmt.Println(loadOne[[]ProductCategorySubcategory](bctx, "ProductCategorySubcategoryListByProductID", fmt.Sprint(p.ID)))
}

// What is a better kv naming? we need to prefix it, e.g. Product.ID since we are passing it down multiple layers.
type KV map[string]any

type Product struct {
	ID    int
	Name  string
	Price int
}

type ProductCategory struct {
	ID        int
	Name      string
	ProductID int
}

type ProductCategorySubcategory struct {
	ID                int
	ProductCategoryID int
	Name              string
}

func buildProduct(kvs ...KV) *Product {
	kv := mergeMaps(kvs...)
	p := &Product{
		ID:    0,
		Name:  "chair",
		Price: 10,
	}
	t := reflect.ValueOf(p).Elem()
	for k, v := range kv {
		f := t.FieldByName(k)
		if !f.IsValid() {
			log.Println("invalid product variant: " + k)
			continue
		}
		f.Set(reflect.ValueOf(v))
	}

	return p
}

func createProduct(ctx context.Context, kvs ...KV) *Product {
	kv := mergeMaps(kvs...)
	p := buildProduct(kv)
	p.ID = 123
	fmt.Println("creating product", p)
	// AFTER
	for k, v := range kv {
		switch k {
		case "WithSubcategories":
			ckv := kv["Category"].(KV)
			ckv = maps.Clone(ckv)

			pkv := maps.Clone(kv)
			delete(pkv, "Category")

			ckv["ProductID"] = p.ID
			ckv["Product"] = pkv

			pc := createProductCategory(ctx, ckv)
			n := v.(int)

			subcategories := make([]ProductCategorySubcategory, n)
			for i := 0; i < n; i++ {
				subcategories[i] = *createProductCategorySubcategory(ctx, append(kvs, map[string]any{
					"ProductCategoryID": pc.ID,
					"Name":              fmt.Sprintf("subcategory_%d", i),
				})...)
			}
			store(ctx, "ProductCategorySubcategoryListByProductID", fmt.Sprint(p.ID), subcategories)
		}
	}

	store(ctx, "Product", fmt.Sprint(p.ID), p)
	return p
}

func buildProductCategory(kvs ...KV) *ProductCategory {
	kv := mergeMaps(kvs...)
	p := &ProductCategory{
		ID:        0,
		Name:      "furniture",
		ProductID: 0,
	}

	t := reflect.ValueOf(p).Elem()
	for k, v := range kv {
		f := t.FieldByName(k)
		if !f.IsValid() {
			log.Println("invalid product category variant: " + k)
			continue
		}
		f.Set(reflect.ValueOf(v))
	}

	return p
}

func createProductCategory(ctx context.Context, kvs ...KV) *ProductCategory {
	kv := mergeMaps(kvs...)
	p := buildProductCategory(kv)
	if p.ProductID == 0 {
		pkv, _ := kv["Product"].(KV)
		p.ProductID = createProduct(ctx, pkv).ID
	}
	p.ID = 65
	fmt.Println("creating product category", p)
	store(ctx, "ProductCategory", fmt.Sprint(p.ID), p)
	return p
}

func buildProductCategorySubcategory(kvs ...KV) *ProductCategorySubcategory {
	kv := mergeMaps(kvs...)
	p := &ProductCategorySubcategory{
		ID:                0,
		Name:              "bed and tables",
		ProductCategoryID: 0,
	}

	t := reflect.ValueOf(p).Elem()
	for k, v := range kv {
		f := t.FieldByName(k)
		if !f.IsValid() {
			fmt.Println("invalid field name", k)
			continue
		}
		f.Set(reflect.ValueOf(v))
	}

	return p
}

func createProductCategorySubcategory(ctx context.Context, kvs ...KV) *ProductCategorySubcategory {
	p := buildProductCategorySubcategory(kvs...)
	if p.ProductCategoryID == 0 {
		p.ProductCategoryID = createProductCategory(ctx, kvs...).ID
	}
	p.ID = 23
	fmt.Println("creating product category subcategory", p)
	store(ctx, "ProductCategorySubcategory", fmt.Sprint(p.ID), p)
	return p
}

// Don't do this, create it under createProduct instead.
func createProductCategorySubcategoryList() {}

type contextKey string

var factoryContextKey contextKey = "factory"

func withGeneratedID(ctx context.Context) context.Context {
	_, ok := generatedID(ctx)
	if ok {
		return ctx
	}
	return context.WithValue(ctx, factoryContextKey, fmt.Sprint(time.Now()))
}

func generatedID(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(factoryContextKey).(string)
	return s, ok
}

func mergeMaps[K comparable, V any, T ~map[K]V](kvs ...T) T {
	if len(kvs) == 1 {
		return kvs[0]
	}
	res := make(map[K]V)
	for _, kv := range kvs {
		for k, v := range kv {
			res[k] = v
		}
	}
	return res
}
```


## Option 3: The simplest

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	ctx := context.Background()
	fmt.Println(createProduct(ctx, "product.expensive"))
	fmt.Println()
	fmt.Println(createProductCategory(ctx))
	fmt.Println()
	fmt.Println(createProductCategory(ctx, "product_category.product_id:111"))
	fmt.Println()
	fmt.Println(createProductCategorySubcategory(ctx))

	agg := createProductAggregate(ctx)
	fmt.Println(agg)
}

type Product struct {
	ID    int
	Name  string
	Price int
}

type ProductCategory struct {
	ID        int
	Name      string
	ProductID int
}

type ProductCategorySubcategory struct {
	ID                int
	ProductCategoryID int
	Name              string
}

func buildProduct(variants ...string) *Product {
	p := &Product{
		ID:    0,
		Name:  "chair",
		Price: 10,
	}
	for _, v := range variants {
		vv := strings.TrimPrefix(v, "product.")
		if vv == v {
			continue
		}
		f, val, _ := strings.Cut(vv, ":")
		switch f {
		case "expensive":
			_ = val
			p.Price = 999_999_999
		default:
			panic("invalid product category variant: " + v)
		}
	}

	return p
}

func createProduct(ctx context.Context, variants ...string) *Product {
	p := buildProduct(variants...)
	p.ID = 123
	fmt.Println("creating product", p)
	return p
}

func buildProductCategory(variants ...string) *ProductCategory {
	p := &ProductCategory{
		ID:        0,
		Name:      "furniture",
		ProductID: 0,
	}

	for _, v := range variants {
		vv := strings.TrimPrefix(v, "product_category.")
		if vv == v {
			continue
		}
		f, val, _ := strings.Cut(vv, ":")
		switch f {
		case "product_id":
			p.ProductID = toInt(val)
		default:
			panic("invalid product category variant: " + v)
		}
	}

	return p
}

func createProductCategory(ctx context.Context, variants ...string) *ProductCategory {
	p := buildProductCategory(variants...)
	if p.ProductID == 0 {
		p.ProductID = createProduct(ctx, variants...).ID
	}
	p.ID = 65
	fmt.Println("creating product category", p)
	return p
}

func buildProductCategorySubcategory(variants ...string) *ProductCategorySubcategory {
	p := &ProductCategorySubcategory{
		ID:                0,
		Name:              "bed and tables",
		ProductCategoryID: 0,
	}

	for _, v := range variants {
		vv := strings.TrimPrefix(v, "product_category_subcategory.")
		if vv == v {
			continue
		}
		f, val, _ := strings.Cut(vv, ":")
		switch f {
		case "product_category_id":
			p.ProductCategoryID = toInt(val)
		default:
			panic("invalid product category variant: " + v)
		}
	}
	return p
}

func createProductCategorySubcategory(ctx context.Context, variants ...string) *ProductCategorySubcategory {
	p := buildProductCategorySubcategory(variants...)
	if p.ProductCategoryID == 0 {
		p.ProductCategoryID = createProductCategory(ctx, variants...).ID
	}
	p.ID = 23
	fmt.Println("creating product category subcategory", p)
	return p
}

// Don't do this, create it under createProductAggregate instead.
func createProductCategorySubcategoryList() {}

type productAggregate struct {
	product              *Product
	productCategory      *ProductCategory
	productSubcategories []ProductCategorySubcategory
}

func createProductAggregate(ctx context.Context, variants ...string) *productAggregate {
	p := createProduct(ctx, variants...)
	pc := createProductCategory(ctx, append(variants, fmt.Sprintf("product_category.product_id:%d", p.ID))...)

	n := 10
	ps := make([]ProductCategorySubcategory, n)
	variants = append(variants, fmt.Sprintf("product_category_subcategory.product_category_id:%d", pc.ID))
	for i := 0; i < n; i++ {
		ps[i] = *createProductCategorySubcategory(ctx, variants...)
	}
	return &productAggregate{
		product:              p,
		productCategory:      pc,
		productSubcategories: ps,
	}
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}
```
