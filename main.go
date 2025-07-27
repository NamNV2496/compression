package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DataDog/zstd"
	"github.com/andybalholm/brotli"
)

const (
	// filePath = "./saokearibank.csv" // 4.4MB
	// filePath = "./body.txt" // 2kB
	filePath = "./100kb.txt" // 107kB
	// how hard should the compressor work to reduce size?
	COMPRESS_LEVEL = 5
)

type CompressRequest struct {
	Data string `json:"data"`
}

type CompressResponse struct {
	Algo           string  `json:"algo"`
	InputSize      int     `json:"input_size"`
	OutputSize     int     `json:"output_size"`
	Ratio          float64 `json:"ratio"`
	CompressRatio  string  `json:"compress_ratio"`
	DurationMillis int64   `json:"duration_ms"`
}

// ad verify input and output
func computeMD5(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func compressGzip(data []byte) (int, int64, error) {
	md5Before := computeMD5(data)
	var buf bytes.Buffer
	start := time.Now()
	w, _ := gzip.NewWriterLevel(&buf, COMPRESS_LEVEL)
	_, err := w.Write(data)
	if err != nil {
		return 0, 0, err
	}
	w.Close()
	output := len(buf.Bytes())
	// Decompress to validate
	gr, err := gzip.NewReader(&buf)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to unzip for validation: %w", err)
	}
	defer gr.Close()

	decompressedData, err := io.ReadAll(gr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	md5After := computeMD5(decompressedData)
	if md5Before != md5After {
		return 0, 0, fmt.Errorf("data mismatch after decompression")
	}
	log.Printf("md5Before: %s, md5After: %s", md5Before, md5After)
	return output, time.Since(start).Milliseconds(), nil
}

func compressBrotli(data []byte) (int, int64, error) {
	var buf bytes.Buffer
	start := time.Now()
	w := brotli.NewWriterLevel(&buf, COMPRESS_LEVEL)
	_, err := w.Write(data)
	if err != nil {
		return 0, 0, err
	}
	w.Close()
	return len(buf.Bytes()), time.Since(start).Milliseconds(), nil
}

func compressZstd(data []byte) (int, int64, error) {
	var buf []byte
	start := time.Now()
	out, err := zstd.CompressLevel(buf, data, COMPRESS_LEVEL)
	if err != nil {
		return 0, 0, err
	}
	return len(out), time.Since(start).Milliseconds(), nil
}

func compressHandler(algo string, compressFunc func([]byte) (int, int64, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Cannot read file", http.StatusInternalServerError)
			return
		}

		outSize, ms, err := compressFunc(data)
		if err != nil {
			http.Error(w, "Compression error", http.StatusInternalServerError)
			return
		}
		ratio := float64(len(data)) / float64(outSize)
		res := CompressResponse{
			Algo:           algo,
			InputSize:      len(data),
			OutputSize:     outSize,
			Ratio:          ratio,
			CompressRatio:  fmt.Sprintf("%.2f%%", float64(outSize)/float64(len(data))*100),
			DurationMillis: ms,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func main() {
	http.HandleFunc("/compress/gzip", compressHandler("gzip", compressGzip))
	http.HandleFunc("/compress/br", compressHandler("brotli", compressBrotli))
	http.HandleFunc("/compress/zstd", compressHandler("zstd", compressZstd))
	log.Println("Listening on :8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// http://localhost:8080/compress/gzip
// http://localhost:8080/compress/br
// http://localhost:8080/compress/zstd
