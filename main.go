package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

var (
	filePath = "./saokearibank.csv"
)

type CompressRequest struct {
	Data string `json:"data"`
}

type CompressResponse struct {
	Algo           string  `json:"algo"`
	InputSize      int     `json:"input_size"`
	OutputSize     int     `json:"output_size"`
	Ratio          float64 `json:"ratio"`
	DurationMillis int64   `json:"duration_ms"`
}

func compressGzip(data []byte) ([]byte, int64, error) {
	var buf bytes.Buffer
	start := time.Now()
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, 0, err
	}
	w.Close()
	return buf.Bytes(), time.Since(start).Milliseconds(), nil
}

func compressBrotli(data []byte) ([]byte, int64, error) {
	var buf bytes.Buffer
	start := time.Now()
	w := brotli.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, 0, err
	}
	w.Close()
	return buf.Bytes(), time.Since(start).Milliseconds(), nil
}

func compressZstd(data []byte) ([]byte, int64, error) {
	var buf bytes.Buffer
	start := time.Now()
	enc, err := zstd.NewWriter(&buf)
	if err != nil {
		return nil, 0, err
	}
	_, err = enc.Write(data)
	if err != nil {
		return nil, 0, err
	}
	enc.Close()
	return buf.Bytes(), time.Since(start).Milliseconds(), nil
}

func compressHandler(algo string, compressFunc func([]byte) ([]byte, int64, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Cannot read file", http.StatusInternalServerError)
			return
		}
		out, ms, err := compressFunc(data)
		if err != nil {
			http.Error(w, "Compression error", http.StatusInternalServerError)
			return
		}

		res := CompressResponse{
			Algo:           algo,
			InputSize:      len(data),
			OutputSize:     len(out),
			Ratio:          float64(len(out)) / float64(len(data)),
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
