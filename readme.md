# Compression

```txt
http://localhost:8080/compress/zip
http://localhost:8080/compress/br
http://localhost:8080/compress/zstd
```

# Benchmark

```
go test -bench=.

BenchmarkCompressGzip-10              39          28772779 ns/op
BenchmarkCompressBrotli-10            46          23064164 ns/op
BenchmarkCompressZstd-10             139           8561519 ns/op

ns/op: time to execute 1 operator
```

```txt
localhost:8080/zstd
localhost:8080/normal
```
