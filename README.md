# Shred function in Go

The current test coverage is 84.3%.
To reach the 100% coverage some tests could be written to check the *shred* package internal functions `overwriteInChunks()` and `writeChunk()`.

Coverage verified using:

```
go test -coverprofile="test.coverage"
go tool cover -html="test.coverage"
```
