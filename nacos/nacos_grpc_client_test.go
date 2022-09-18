package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var grpcClientTest = NewNacosGrpcClientTest()

func NewNacosGrpcClientTest() *NacosGrpcClient {
	grpcClient, err := NewNacosGrpcClient("", []string{"106.52.77.111:8848"}, nacosClientTest)
	if err != nil {
		fmt.Println("init grpc client failed")
	}
	return grpcClient
}

func TestGetAllServicesInfo(t *testing.T) {
	services := grpcClientTest.GetAllServicesInfo()
	if len(services) > 0 {
		t.Log("GrpcClient get all servicesInfo passed")
	} else {
		t.Log("GrpcClient get all servicesInfo empty")
	}
}

func TestGetService(t *testing.T) {
	services := grpcClientTest.GetAllServicesInfo()
	serviceMap := NewConcurrentMap()
	for _, serviceName := range services {
		service := grpcClientTest.GetService(serviceName)
		if assert.NotNil(t, service) {
			serviceMap.Set(serviceName, service)
		}
	}
	if serviceMap.Count() == len(services) {
		t.Log("GrpcClient get service passed")
	} else {
		t.Error("GrpcClient get service error")
	}
}

func TestSubscribe(t *testing.T) {
	doms := grpcClientTest.GetAllServicesInfo()
	for _, dom := range doms {
		err := grpcClientTest.Subscribe(dom)
		if err != nil {
			t.Error("GrpcClient subscribe service error")
			return
		}
	}
	t.Log("GrpcClient subscribe service passed")
}

func TestCallback(t *testing.T) {
	services := model.Service{
		Name:        "DEFAULT_GROUP@@demo.go",
		CacheMillis: 1000,
		Hosts: []model.Instance{
			{
				InstanceId:  "10.10.10.10-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.10",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      true,
				Healthy:     true,
			},
			{
				InstanceId:  "10.10.10.11-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.11",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      true,
				Healthy:     true,
			},
			{
				InstanceId:  "10.10.10.12-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.12",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      true,
				Healthy:     false,
			},
			{
				InstanceId:  "10.10.10.13-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.13",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      false,
				Healthy:     true,
			},
			{
				InstanceId:  "10.10.10.14-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.14",
				Weight:      0,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      true,
				Healthy:     true,
			},
		},
		Checksum:    "3bbcf6dd1175203a8afdade0e77a27cd1528787794594",
		LastRefTime: 1528787794594, Clusters: "a"}

	grpcClientTest.nacosClient.GetDomainCache().Set("demo.go", services)

	newServices := model.Service{
		Name:        "DEFAULT_GROUP@@demo.go",
		CacheMillis: 1000,
		Hosts: []model.Instance{
			{
				InstanceId:  "10.10.10.10-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.10",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      true,
				Healthy:     true,
			},
			{
				InstanceId:  "10.10.10.11-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.11",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      true,
				Healthy:     true,
			},
			{
				InstanceId:  "10.10.10.12-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.12",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      true,
				Healthy:     false,
			},
			{
				InstanceId:  "10.10.10.13-80-a-DEMO",
				Port:        80,
				Ip:          "10.10.10.13",
				Weight:      1,
				Metadata:    map[string]string{},
				ClusterName: "a",
				ServiceName: "DEFAULT_GROUP@@demo.go",
				Enable:      false,
				Healthy:     true,
			},
		},
		Checksum:    "3bbcf6dd1175203a8afdade0e77a27cd1528787794594",
		LastRefTime: 1528787794594, Clusters: "a"}
	grpcClientTest.Callback(newServices.Hosts, nil)

	s, _ := grpcClientTest.nacosClient.GetDomainCache().Get("demo.go")

	updateServices := s.(model.Service)

	if len(newServices.Hosts) == len(updateServices.Hosts) {
		t.Log("GrpcClient Service SubscribeCallback passed")
	} else {
		t.Error("GrpcClient Service SubscribeCallback error")
	}
}
