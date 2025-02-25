package k8sresources

import (
	"fmt"
	"strings"

	"github.com/bonnefoa/kubectl-fzf/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

const ServiceHeader = "Namespace Name Type ClusterIp Ports Selector Age Labels\n"

// Service is the summary of a kubernetes service
type Service struct {
	ResourceMeta
	serviceType string
	clusterIP   string
	ports       []string
	selectors   []string
}

// NewServiceFromRuntime builds a pod from informer result
func NewServiceFromRuntime(obj interface{}, config CtorConfig) K8sResource {
	s := &Service{}
	s.FromRuntime(obj, config)
	return s
}

// FromRuntime builds object from the informer's result
func (s *Service) FromRuntime(obj interface{}, config CtorConfig) {
	service := obj.(*corev1.Service)
	s.FromObjectMeta(service.ObjectMeta)
	s.serviceType = string(service.Spec.Type)
	s.clusterIP = service.Spec.ClusterIP
	s.ports = make([]string, len(service.Spec.Ports))
	for k, v := range service.Spec.Ports {
		if v.NodePort > 0 {
			s.ports[k] = fmt.Sprintf("%s:%d/%d", v.Name, v.Port, v.NodePort)
		} else {
			s.ports[k] = fmt.Sprintf("%s:%d", v.Name, v.Port)
		}
	}
	s.selectors = util.JoinStringMap(service.Spec.Selector, ExcludedLabels, "=")
}

// HasChanged returns true if the resource's dump needs to be updated
func (s *Service) HasChanged(k K8sResource) bool {
	oldService := k.(*Service)
	return (util.StringSlicesEqual(s.ports, oldService.ports) ||
		util.StringSlicesEqual(s.selectors, oldService.selectors) ||
		util.StringMapsEqual(s.labels, oldService.labels))
}

// ToString serializes the object to strings
func (s *Service) ToString() string {
	portList := util.JoinSlicesOrNone(s.ports, ",")
	selectorList := util.JoinSlicesOrNone(s.selectors, ",")
	line := strings.Join([]string{s.namespace,
		s.name,
		s.serviceType,
		s.clusterIP,
		portList,
		selectorList,
		s.resourceAge(),
		s.labelsString(),
	}, " ")
	return fmt.Sprintf("%s\n", line)
}
