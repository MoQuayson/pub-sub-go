package broker

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"testing"
)

func TestGrpcBrokerServerWithInMemoryStorage(t *testing.T) {
	brokerSrv := broker.NewBroker(&models.BrokerConfig{
		Host:      "0.0.0.0",
		Port:      "7000",
		Storage:   models.InMemoryStorageType,
		Transport: models.GrpcTransport,
	})

	if err := brokerSrv.Start(); err != nil {
		t.Error(err)
	}

	t.Log("test done")
}

func TestGrpcBrokerServerWithDiskStorage(t *testing.T) {
	brokerSrv := broker.NewBroker(&models.BrokerConfig{
		Host:            "0.0.0.0",
		Port:            "50051",
		Storage:         models.DiskStorageType,
		StorageLocation: utils.ConvertToPointerString(logDir),
		Transport:       models.GrpcTransport,
	})

	if err := brokerSrv.Start(); err != nil {
		t.Error(err)
	}

	t.Log("test done")
}
