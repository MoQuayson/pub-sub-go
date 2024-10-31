package enums

type StorageType string

const (
	StorageType_InMemory StorageType = "InMemory"
	StorageType_Redis    StorageType = "Redis"
	StorageType_Disk     StorageType = "Disk"
)
