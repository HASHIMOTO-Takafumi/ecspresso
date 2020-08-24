package appspec

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"gopkg.in/yaml.v2"
)

var (
	Version    = "0.0"
	TargetType = "AWS::ECS::Service"
)

type AppSpec struct {
	Version   *string     `yaml:"version"`
	Resources []*Resource `yaml:"Resources,omitempty"`
	Hooks     []*Hook     `yaml:"Hooks,omitempty"`
}

func New() *AppSpec {
	return &AppSpec{
		Version: &Version,
	}
}

func (a *AppSpec) String() string {
	b, _ := yaml.Marshal(a)
	return string(b)
}

func NewWithService(sv *ecs.Service, tdArn string) (*AppSpec, error) {
	if len(sv.LoadBalancers) == 0 {
		return nil, errors.New("require LoadBalancers")
	}
	spec := New()
	resource := &Resource{
		TargetService: &TargetService{
			Type: aws.String(TargetType),
			Properties: &Properties{
				TaskDefinition: aws.String(tdArn),
				LoadBalancerInfo: &LoadBalancerInfo{
					ContainerName: sv.LoadBalancers[0].ContainerName,
					ContainerPort: sv.LoadBalancers[0].ContainerPort,
				},
				PlatformVersion: sv.PlatformVersion,
			},
		},
	}
	if sv.NetworkConfiguration != nil && sv.NetworkConfiguration.AwsvpcConfiguration != nil {
		cfg := sv.NetworkConfiguration.AwsvpcConfiguration
		resource.TargetService.Properties.NetworkConfiguration = &NetworkConfiguration{
			AwsvpcConfiguration: &AwsVpcConfiguration{
				Subnets:        cfg.Subnets,
				SecurityGroups: cfg.SecurityGroups,
				AssignPublicIp: cfg.AssignPublicIp,
			},
		}
	}
	spec.Resources = append(spec.Resources, resource)
	return spec, nil
}

type Resource struct {
	TargetService *TargetService `yaml:"TargetService,omitempty"`
}

type TargetService struct {
	Type       *string     `yaml:"Type,omitempty"`
	Properties *Properties `yaml:"Properties,omitempty"`
}

type Properties struct {
	TaskDefinition       *string               `yaml:"TaskDefinition,omitempty"`
	LoadBalancerInfo     *LoadBalancerInfo     `yaml:"LoadBalancerInfo,omitempty"`
	PlatformVersion      *string               `yaml:"PlatformVersion,omitempty"`
	NetworkConfiguration *NetworkConfiguration `yaml:"NetworkConfiguration,omitempty"`
}

type LoadBalancerInfo struct {
	ContainerName *string `yaml:"ContainerName"`
	ContainerPort *int64  `yaml:"ContainerPort"`
}

type NetworkConfiguration struct {
	AwsvpcConfiguration *AwsVpcConfiguration `yaml:"AwsvpcConfiguration,omitempty"`
}

type AwsVpcConfiguration struct {
	AssignPublicIp *string   `yaml:"AssignPublicIp,omitempty"`
	SecurityGroups []*string `yaml:"SecurityGroups,omitempty"`
	Subnets        []*string `yaml:"Subnets,omitempty"`
}

type Hook struct {
	BeforeInstall         string `yaml:"BeforeInstall,omitempty"`
	AfterInstall          string `yaml:"AfterInstall,omitempty"`
	AfterAllowTestTraffic string `yaml:"AfterAllowTestTraffic,omitempty"`
	BeforeAllowTraffic    string `yaml:"BeforeAllowTraffic,omitempty"`
	AfterAllowTraffic     string `yaml:"AfterAllowTraffic,omitempty"`
}
