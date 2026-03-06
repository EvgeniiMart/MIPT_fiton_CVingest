package storage

func NewMockStorage() *Storage {

	return &Storage{
		Seaweed:  &SeaweedMock{},
		Postgres: &PostgresMock{},
		Redis:    &RedisMock{},
	}
}
