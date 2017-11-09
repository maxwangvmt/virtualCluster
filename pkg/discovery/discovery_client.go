package discovery

import (
	"fmt"
	"github.com/golang/glog"

	"virtualCluster/pkg/registration"
	"virtualCluster/pkg/target"

	sdkprobe "github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type DiscoveryClient struct {
	targetConfig *TargetConf
	cluster      *target.ClusterHandler
}

func NewDiscoveryClient(targetConfig *TargetConf, handler *target.ClusterHandler) *DiscoveryClient {
	return &DiscoveryClient{
		targetConfig: targetConfig,
		cluster:      handler,
	}
}

func (dc *DiscoveryClient) String() string {
	return fmt.Sprintf("%+v\n%v", dc.targetConfig, dc.cluster.String())
}

func (dc *DiscoveryClient) GetAccountValues() *sdkprobe.TurboTargetInfo {
	var accountValues []*proto.AccountValue

	targetConf := dc.targetConfig
	// Convert all parameters in clientConf to AccountValue list
	targetID := registration.TargetIdentifierField
	accVal := &proto.AccountValue{
		Key:         &targetID,
		StringValue: &targetConf.Address,
	}
	accountValues = append(accountValues, accVal)

	username := registration.Username
	accVal = &proto.AccountValue{
		Key:         &username,
		StringValue: &targetConf.Username,
	}
	accountValues = append(accountValues, accVal)

	password := registration.Password
	accVal = &proto.AccountValue{
		Key:         &password,
		StringValue: &targetConf.Password,
	}
	accountValues = append(accountValues, accVal)

	targetInfo := sdkprobe.NewTurboTargetInfoBuilder(targetConf.ProbeCategory, targetConf.TargetType, targetID, accountValues).Create()

	glog.V(2).Infof("Got AccountValues for target:%v", targetConf.Address)
	return targetInfo
}

func (dc *DiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.V(2).Infof("begin to validating target...")

	return &proto.ValidationResponse{}, nil
}

func printDTOs(dtos []*proto.EntityDTO) string {
	msg := ""
	for _, dto := range dtos {
		line := fmt.Sprintf("%+v", dto)
		msg = msg + "\n" + line
	}

	return msg
}

func (dc *DiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	glog.V(2).Infof("begin to discovery target...")

	resultDTOs, err := dc.cluster.GenerateClusterDTOs()
	if err != nil {
		glog.Errorf("failed to generate DTOs: %v", err)
		resultDTOs = []*proto.EntityDTO{}
	}

	glog.V(2).Infof("end of discoverying target. [%d]", len(resultDTOs))
	glog.V(4).Infof("DTOs:\n%s", printDTOs(resultDTOs))

	response := &proto.DiscoveryResponse{
		EntityDTO: resultDTOs,
	}

	return response, nil
}
