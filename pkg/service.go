package pkg

import (
	"errors"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"os/exec"
	"strconv"
	"sync"
)

// SerfPublisherInterface for the SerfPublisher
type SerfPublisherInterface interface {
	Publish(service v1.Service) (v1.Service, error)
	Unpublish(key string) (v1.Service, error)
}

// SerfPublisher is simple annotator service.
type SerfPublisher struct {
	client kubernetes.Interface
	logger Logger
	reg    sync.Map
}

// NewSerfPublisher returns a new SerfPublisher.
func NewSerfPublisher(k8sCli kubernetes.Interface, logger Logger) *SerfPublisher {
	return &SerfPublisher{
		client: k8sCli,
		logger: logger,
		reg:    sync.Map{},
	}
}

// Publish will add a new service trhought serf
func (s *SerfPublisher) Publish(service v1.Service) (v1.Service, error) {
	newService := service.DeepCopy()
	cmd := exec.Command("/usr/sbin/avahi-ps", "publish", newService.ObjectMeta.Name, "kubernetes", strconv.Itoa(int(newService.Spec.Ports[0].NodePort)), string(newService.Spec.Type))
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Infof("cmd.Run() failed with %s\n", err)
	}
	s.logger.Infof("command \n%s\n", out)
	s.reg.Store(newService.Name, newService)
	return *newService, nil
}

// Publish will add a new service trhought serf
func (s *SerfPublisher) Unpublish(key string) (v1.Service, error) {
	_, serviceName, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return v1.Service{}, errors.New("can't split namespace and name")
	}
	svc, ok := s.reg.Load(serviceName)
	if !ok {
		return v1.Service{}, errors.New("unknown deleted object")
	}
	newService := svc.(*v1.Service)
	cmd := exec.Command("/usr/sbin/avahi-ps", "unpublish", "kubernetes", strconv.Itoa(int(newService.Spec.Ports[0].NodePort)))
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Infof("cmd.Run() failed with %s\n", err)
	}
	s.logger.Infof("command \n%s\n", out)
	return *newService, nil
}