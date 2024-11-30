package models

type StorageType string

const (
	InMemoryStorageType StorageType = "InMemory"
	RedisStorageType    StorageType = "Redis"
	DiskStorageType     StorageType = "Disk"
)
