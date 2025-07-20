package main

// import (
// 	"os"
// 	"testing"
// )

// var benchData []byte

// func init() {
// 	// Đọc data 1 lần khi benchmark bắt đầu
// 	var err error
// 	benchData, err = os.ReadFile(filePath)
// 	if err != nil {
// 		panic("Cannot read file for benchmark: " + err.Error())
// 	}
// }

// func BenchmarkCompressGzip(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, _, err := compressGzip(benchData)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// }

// func BenchmarkCompressBrotli(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, _, err := compressBrotli(benchData)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// }

// func BenchmarkCompressZstd(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, _, err := compressZstd(benchData)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 	}
// }

// go test -bench=.
