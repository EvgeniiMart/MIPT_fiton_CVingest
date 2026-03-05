package main

import (
	"fmt"
	"net/http"

	"github.com/EvgeniiMart/MIPT_fiton_CVingest/internal/receive"
)

const schemaPath = "data/json_schemas/cv_protocol.schema.json"

func main() {

	http.HandleFunc("/api/v1/ingest/batch",
		receive.IngestBatchHandler(schemaPath))

	fmt.Println("Server started on :8080")

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
