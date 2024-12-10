package broker

import (
	"github.com/MoQuayson/pub-sub-go/pkg/broker"
	"github.com/MoQuayson/pub-sub-go/pkg/utils"
	"github.com/MoQuayson/pub-sub-go/pkg/utils/models"
	"testing"
)

const logDir = "../../bin/logs"

func TestRpcBrokerServerWithInMemoryStorage(t *testing.T) {
	brokerSrv := broker.NewBroker(&models.BrokerConfig{
		Host:      "0.0.0.0",
		Port:      "7000",
		Storage:   models.InMemoryStorageType,
		Transport: models.DefaultTransport,
	})

	if err := brokerSrv.Start(); err != nil {
		t.Error(err)
	}

	t.Log("test done")
}

func TestRpcBrokerServerWithDiskStorage(t *testing.T) {
	brokerSrv := broker.NewBroker(&models.BrokerConfig{
		Host:            "0.0.0.0",
		Port:            "7000",
		Storage:         models.DiskStorageType,
		StorageLocation: utils.ConvertToPointerString(logDir),
		Transport:       models.DefaultTransport,
	})

	if err := brokerSrv.Start(); err != nil {
		t.Error(err)
	}

	t.Log("test done")
}
