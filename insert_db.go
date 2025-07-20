package main

// import (
// 	"database/sql"
// 	"encoding/csv"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/DataDog/zstd"
// 	_ "github.com/lib/pq"
// )

// const (
// 	filePath = "./saokearibank.csv" // 4.4MB
// 	// how hard should the compressor work to reduce size?
// 	COMPRESS_LEVEL = 1
// )

// func compressHandler(compressFunc func(string) []byte) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		file, err := os.Open(filePath)
// 		if err != nil {
// 			http.Error(w, "Cannot open file", http.StatusInternalServerError)
// 			return
// 		}
// 		defer file.Close()

// 		reader := csv.NewReader(file)
// 		records, err := reader.ReadAll()
// 		if err != nil {
// 			http.Error(w, "Cannot read CSV", http.StatusInternalServerError)
// 			return
// 		}

// 		connStr := "user=root password=root dbname=zstd sslmode=disable"
// 		db, err := sql.Open("postgres", connStr)
// 		if err != nil {
// 			http.Error(w, "Cannot connect to DB", http.StatusInternalServerError)
// 			return
// 		}
// 		defer db.Close()

// 		for i, record := range records {
// 			if i == 0 {
// 				continue // Skip header
// 			}
// 			_, err := db.Exec("INSERT INTO zstd (date, remark, credit, balance, ref) VALUES ($1, $2, $3, $4, $5)",
// 				record[0],
// 				compressFunc(record[1]),
// 				record[2],
// 				record[3],
// 				record[4])
// 			if err != nil {
// 				http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 		take := fmt.Sprintf("CSV data inserted successfully! take: %d ms", time.Since(start).Milliseconds())
// 		if _, err := w.Write([]byte(take)); err != nil {
// 			log.Printf("failed to write response: %v", err)
// 		}
// 	}
// }

// func compressZstd(data string) []byte {
// 	var buf []byte
// 	out, err := zstd.CompressLevel(buf, []byte(data), COMPRESS_LEVEL)
// 	if err != nil {
// 		return nil
// 	}
// 	return out
// }
// func main() {

// 	http.HandleFunc("/compress/zstd", compressHandler(compressZstd))
// 	log.Println("Listening on :8080 ...")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// // curl -L "http://localhost:8080/compress/zstd"
