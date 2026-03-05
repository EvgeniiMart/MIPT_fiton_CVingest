package receive

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AcceptResponse struct {
	BatchID string `json:"batch_id"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

func IngestBatchHandler(schemaPath string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "used method not allowed, should be POST",
				http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			writeError(w, "invalid multipart")
			return
		}

		jsonData, err := ParseMultipart(r.MultipartForm)
		if err != nil {
			writeError(w, err.Error())
			return
		}

		err = ValidateJSON(schemaPath, string(jsonData))
		if err != nil {
			writeError(w, err.Error())
			return
		}

		var parsed struct {
			Metadata struct {
				BatchID string `json:"batch_id"`
			} `json:"metadata"`
		}

		if err := json.Unmarshal(jsonData, &parsed); err != nil {
			writeError(w, "parsing error on the receiver side: "+
				err.Error())
			return
		}

		resp := AcceptResponse{
			BatchID: parsed.Metadata.BatchID,
			Status:  "accepted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(resp)

		fmt.Printf("Batch accepted: %s\n", parsed.Metadata.BatchID)
	}
}
