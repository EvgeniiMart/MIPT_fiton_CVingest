package storage

import (
	"log"

	"github.com/google/uuid"
)

type SeaweedMock struct{}

func (s *SeaweedMock) Upload(data []byte) (string, error) {
	fid := uuid.New().String()

	log.Printf("[SeaweedFSMock] upload size=%d -> fid=%s", len(data), fid)

	return fid, nil
}
