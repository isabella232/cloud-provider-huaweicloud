/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package huaweicloud

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

const (
	ELBIDAnnotation            = "kubernetes.io/elb.id"
	ELBClassAnnotation         = "kubernetes.io/elb.class"
	ELBMarkAnnotation          = "kubernetes.io/elb.mark"
	VPCIDAnnotation            = "kubernetes.io/elb.vpc.id"
	ELBSessionAffinityMode     = "kubernetes.io/session-affinity-mode"
	ELBSessionSourceIP         = "SOURCE_IP"
	Ping                       = "ping"
	Pong                       = "pong"
	HealthzCCE                 = "cce-healthz"
	ListenerDescription        = "Attention! It is auto-generated by CCE service, do not modify!"
	DefaultSessionAffinityTime = 1440
)

type LoadBalancerOpts struct {
	Apiserver    string       `json:"apiserver"`
	SecretName   string       `json:"secretName"`
	SignerType   string       `json:"signerType"`
	ELBAlgorithm ELBAlgorithm `json:"elbAlgorithm"`
	TenantId     string       `json:"tenantId"`
	Region       string       `json:"region"`
	VPCId        string       `json:"vpcId"`
	SubnetId     string       `json:"subnetId"`
	ECSEndpoint  string       `json:"ecsEndpoint"`
	ELBEndpoint  string       `json:"elbEndpoint"`
	ALBEndpoint  string       `json:"albEndpoint"`
	NATEndpoint  string       `json:"natEndpoint"`
	VPCEndpoint  string       `json:"vpcEndpoint"`
}

// ELBAlgorithm
type ELBAlgorithm string

const (
	ELBAlgorithmRR  ELBAlgorithm = "roundrobin"
	ELBAlgorithmLC  ELBAlgorithm = "leastconn"
	ELBAlgorithmSRC ELBAlgorithm = "source"
)

type LoadBalanceVersion int

const (
	VersionNotNeedLB LoadBalanceVersion = iota //if the service type is not LoadBalancer
	VersionELB
	VersionALB
	VersionNAT
)

func GetHealthCheckPort(service *v1.Service) *v1.ServicePort {
	for _, port := range service.Spec.Ports {
		if port.Name == HealthzCCE {
			return &port
		}
	}
	return nil
}

func GetListenerName(service *v1.Service) string {
	return string(service.UID)
}

// to suit for old version
// if the elb has been created with the old version
// its listener name is service.name+service.uid
func GetOldListenerName(service *v1.Service) string {
	return strings.Replace(service.Name+"_"+string(service.UID), ".", "_", -1)
}

func GetSessionAffinity(service *v1.Service) bool {
	if service.Annotations[ELBSessionAffinityMode] == ELBSessionSourceIP {
		return true
	}
	return false
}

// if the node not health, it will not be added to ELB
func CheckNodeHealth(node *v1.Node) (bool, error) {
	conditionMap := make(map[v1.NodeConditionType]*v1.NodeCondition)
	for i := range node.Status.Conditions {
		cond := node.Status.Conditions[i]
		conditionMap[cond.Type] = &cond
	}

	status := false
	if condition, ok := conditionMap[v1.NodeReady]; ok {
		if condition.Status == v1.ConditionTrue {
			status = true
		} else {
			status = false
		}
	}

	if node.Spec.Unschedulable {
		status = false
	}

	return status, nil
}