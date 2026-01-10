package ocr3capability

import (
	ocr4 "github.com/smartcontractkit/chainlink-common/pkg/loop/relayer/pluginprovider/ocr2"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/relayer/pluginprovider/ocr3"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	ocr3pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/net"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var (
	_ types.OCR3CapabilityProvider = (*ProviderClient)(nil)
	_ goplugin.GRPCClientConn      = (*ProviderClient)(nil)
)

type ProviderClient struct {
	*ocr4.PluginProviderClient
	ocr3ContractTransmitter ocr3types.ContractTransmitter[[]byte]
}

func (p *ProviderClient) OCR3ContractTransmitter() ocr3types.ContractTransmitter[[]byte] {
	return p.ocr3ContractTransmitter
}

func NewProviderClient(b *net.BrokerExt, cc grpc.ClientConnInterface) *ProviderClient {
	m := &ProviderClient{
		PluginProviderClient:    ocr4.NewPluginProviderClient(b.WithName("OCR3CapabilityProviderClient"), cc),
		ocr3ContractTransmitter: ocr3.NewContractTransmitterClient(b.WithName("OCR3ContractTransmitter"), cc),
	}

	return m
}

type ProviderServer struct{}

func (m ProviderServer) ConnToProvider(conn grpc.ClientConnInterface, broker net.Broker, brokerCfg net.BrokerConfig) types.OCR3CapabilityProvider {
	be := &net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}
	return NewProviderClient(be, conn)
}

func RegisterProviderServices(s *grpc.Server, provider types.OCR3CapabilityProvider) {
	ocr4.RegisterPluginProviderServices(s, provider)
	ocr3pb.RegisterContractTransmitterServer(s, ocr3.NewContractTransmitterServer(provider.OCR3ContractTransmitter()))
}
