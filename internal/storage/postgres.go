package storage

import "log"

type PostgresMock struct{}

func (p *PostgresMock) InsertBatch(
	batchID string,
	cameraID string,
	capturedAt string,
	seaweedFids map[string]string,
) error {

	log.Printf(
		"[PostgresMock] INSERT batches id=%s camera_id=%s "+
			"captured_at=%s fids=%v status=pending",
		batchID,
		cameraID,
		capturedAt,
		seaweedFids,
	)

	return nil
}
