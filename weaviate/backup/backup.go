package backup

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
)

const (
	BACKEND_FILESYSTEM = "filesystem"
	BACKEND_S3         = "s3"
	BACKEND_GCS        = "gcs"
	BACKEND_AZURE      = "azure"
)

type API struct {
	connection *connection.Connection
}

func New(connection *connection.Connection) *API {
	return &API{connection}
}

// Creator creates backup creator builder
func (s *API) Creator() *BackupCreator {
	return &BackupCreator{
		connection:   s.connection,
		statusGetter: s.CreateStatusGetter(),
		compression: Compression{
			CPUPercentage: 50,
			ChunkSize:     128,
			Level:         BackupConfigCompressionLevelDefaultCompression,
		},
	}
}

// CreateStatusGetter creates create status getter builder
func (s *API) CreateStatusGetter() *BackupCreateStatusGetter {
	return &BackupCreateStatusGetter{
		connection: s.connection,
	}
}

// Restorer creates restorer builder
func (s *API) Restorer() *BackupRestorer {
	return &BackupRestorer{
		connection:    s.connection,
		statusGetter:  s.RestoreStatusGetter(),
		cpuPercentage: 50,
	}
}

// RestoreStatusGetter creates restore status getter builder
func (s *API) RestoreStatusGetter() *BackupRestoreStatusGetter {
	return &BackupRestoreStatusGetter{
		connection: s.connection,
	}
}

var (
	// BackupConfigCompressionLevelDefaultCompression captures enum value "DefaultCompression"
	BackupConfigCompressionLevelDefaultCompression string = "DefaultCompression"

	// BackupConfigCompressionLevelBestSpeed captures enum value "BestSpeed"
	BackupConfigCompressionLevelBestSpeed string = "BestSpeed"

	// BackupConfigCompressionLevelBestCompression captures enum value "BestCompression"
	BackupConfigCompressionLevelBestCompression string = "BestCompression"
)

// Compression is the compression configuration.
type Compression struct {
	// Level is one of DefaultCompression, BestSpeed, BestCompression
	Level string

	// ChunkSize represents the desired size for chunks between 1 - 512  MB
	// However, during compression, the chunk size might
	// slightly deviate from this value, being either slightly
	// below or above the specified size
	ChunkSize int

	// CPUPercentage desired CPU core utilization (1%-80%), default: 50%
	CPUPercentage int
}
