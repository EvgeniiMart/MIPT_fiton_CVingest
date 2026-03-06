package storage

import (
	"encoding/json"
	"log"
)

type RedisMock struct{}

func (r *RedisMock) XAdd(
	batchID string,
	metadata json.RawMessage,
	seaweedFids map[string]string,
) error {
	log.Printf(
		"[RedisMock] XADD ingest-queue batch_id=%s fids=%v",
		batchID,
		seaweedFids,
	)

	return nil
}
