package broker

import (
	"github.com/MoQuayson/go-event-bridge/internal/rpc/broker"
	"github.com/MoQuayson/go-event-bridge/tests/shared"
	"testing"
)

func TestBrokerServer(t *testing.T) {
	//config := &models.BrokerConfig{
	//	Host:    "0.0.0.0",
	//	Port:    "12345",
	//	Storage: enums.StorageType_InMemory,
	//}

	brokerSrv := broker.NewRpcBroker(shared.BrokerCfg)
	brokerSrv.Start()

	t.Log("test done")
}
