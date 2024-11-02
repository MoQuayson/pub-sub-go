package broker

import (
	"github.com/MoQuayson/pub-sub-go/internal/rpc/broker/services"
	"github.com/MoQuayson/pub-sub-go/tests/shared"
	"testing"
)

func TestBrokerServer(t *testing.T) {
	brokerSrv := services.NewBrokerService(shared.BrokerCfg)
	brokerSrv.Start()

	t.Log("test done")
}
