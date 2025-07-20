package main

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/DataDog/zstd"
// )

// var (
// 	filePath       = "saokearibank.csv"
// 	COMPRESS_LEVEL = 5
// )

// func decompressZstd(data []byte) ([]byte, int64, error) {
// 	var buf []byte
// 	start := time.Now()
// 	out, err := zstd.Decompress(buf, data)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return out, time.Since(start).Milliseconds(), nil
// }

// func compressZstd(data []byte) ([]byte, int64, error) {
// 	var buf []byte
// 	start := time.Now()
// 	out, err := zstd.CompressLevel(buf, data, COMPRESS_LEVEL)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return out, time.Since(start).Milliseconds(), nil
// }

// func sendNormal(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		http.Error(w, "Cannot read file", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(data)
// }

// func sendCompress(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		http.Error(w, "Cannot read file", http.StatusInternalServerError)
// 		return
// 	}
// 	out, ms, err := compressZstd(data)
// 	if err != nil {
// 		return
// 	}
// 	// resp, _, err := decompressZstd(out)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	log.Printf("compression time: %d ms\n", ms)
// 	w.Write(out)
// }

// func main() {
// 	http.Handle("/normal", http.HandlerFunc(sendNormal))
// 	http.Handle("/zstd", http.HandlerFunc(sendCompress))
// 	http.ListenAndServe(":8080", nil)
// }

// localhost:8080/zstd
// localhost:8080/normal
