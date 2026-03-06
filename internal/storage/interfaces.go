package storage

import "encoding/json"

type Seaweed interface {
	Upload(data []byte) (string, error)
}

type Postgres interface {
	InsertBatch(
		batchID string,
		cameraID string,
		capturedAt string,
		seaweedFids map[string]string,
	) error
}

type Redis interface {
	XAdd(
		batchID string,
		metadata json.RawMessage,
		seaweedFids map[string]string,
	) error
}

type Storage struct {
	Seaweed  Seaweed
	Postgres Postgres
	Redis    Redis
}
