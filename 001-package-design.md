# Package Design


Tips on writing better package

## Start with a clear problem

What is the desired input/output? Can it be done with the standard package instead?

e.g.
- standardize caching interface
- implement distributed locking

## Categorizing package

- helper/utils
- extensions
- common abstraction (your companies core package) for cross cutting concerns (observability)
- standards (your own definition of standard libraries)
- test utils (setting up db migration in docker etc, snapshot testing)

## Methods

Use clear verbs for method names. The mame should clearly indicate what operation it should do. There are some commonly used ones
- `Store/Load` which is preferred over `Get/Set` to indicate the data store is remote
- `Exec` when the signature only returns error
- `Do` when it returns values

Pass context as first argument.

## Sync/dsync

For data structure that can operate concurrently, place them under `sync` package. If the package operates concurrently in distributed settings, place them inder `dsync`.

## Config and Options

Pass a `Config` with public fields in the constructor, and allow passing `Options` when invoking.

This has two benefits
- global configuration (or default)
- local options for overwriting the global config

In the readme, show the construction of the package with all options available.



## Storage

Some packages works with different storage (redis, postgres). Design an interface that allows swapping storage.

The implementation doesnt need to know about the underlying storage.

There are some exceptions, e.g. when it comes to consistency etc like database transaction behaviour. But they should best be documented for clarity if it differs from original design.

## Types

If the package is mostly data structure manipulation or extension of primitives that are not value objects, place them under `types`.

E.g.
- stringcase
- timerange
- interval
- collections
- sets
- orderedmap
- ttlmap

