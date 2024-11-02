package broker

import (
	"github.com/MoQuayson/pub-sub-go/internal/rpc/broker"
	"github.com/MoQuayson/pub-sub-go/tests/shared"
	"testing"
)

func TestBrokerServer(t *testing.T) {
	brokerSrv := broker.NewBrokerService(shared.BrokerCfg)
	brokerSrv.Start()

	t.Log("test done")
}
