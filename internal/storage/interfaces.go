package storage

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

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

func SeaweedUpload(batchField string,
	storage *Storage) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(batchField)
	fid, err := storage.Seaweed.Upload(data)
	if err != nil {
		return "", fmt.Errorf("Error processing batch: %w", err)
	}
	return fid, nil
}

func ProcessBatch(batchJSON []byte, storage *Storage) error {
	var batch map[string]interface{}

	err := json.Unmarshal(batchJSON, &batch)
	if err != nil {
		return fmt.Errorf("Error processing batch: %w", err)
	}

	metadata := batch["metadata"].(map[string]interface{})
	batchID := metadata["batch_id"].(string)
	camera := metadata["camera"].(map[string]interface{})
	cameraType := camera["type"].(string)
	capturedAt := metadata["captured_at"].(string)

	seaweed_fids := map[string]string{}
	if field, ok := batch["image"].(string); ok {
		seaweed_fids["image"], err = SeaweedUpload(field, storage)
		if err != nil {
			return fmt.Errorf("Error processing batch: %w", err)
		}
	}
	if field, ok := batch["mask"].(string); ok {
		seaweed_fids["mask"], err = SeaweedUpload(field, storage)
		if err != nil {
			return fmt.Errorf("Error processing batch: %w", err)
		}
	}
	if field, ok := batch["raw"].(string); ok {
		seaweed_fids["raw"], err = SeaweedUpload(field, storage)
		if err != nil {
			return fmt.Errorf("Error processing batch: %w", err)
		}
	}

	err = storage.Postgres.InsertBatch(batchID, cameraType,
		capturedAt, seaweed_fids)
	if err != nil {
		return fmt.Errorf("Error processing batch: %w", err)
	}

	err = storage.Redis.XAdd(batchID, batchJSON, seaweed_fids)
	if err != nil {
		return fmt.Errorf("Error processing batch: %w", err)
	}

	return nil
}
