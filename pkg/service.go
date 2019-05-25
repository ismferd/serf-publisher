package pkg

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"os/exec"
)

// SerfPublisherInterface for the SerfPublisher
type SerfPublisherInterface interface {
	Publish(service v1.Service) (v1.Service, error)
}

// SerfPublisher is simple annotator service.
type SerfPublisher struct {
	client kubernetes.Interface
	logger Logger
}

// NewSerfPublisher returns a new SerfPublisher.
func NewSerfPublisher(k8sCli kubernetes.Interface, logger Logger) *SerfPublisher {
	return &SerfPublisher{
		client: k8sCli,
		logger: logger,
	}
}

// Publish will add the kube2iam annotation to Deployment objects containing the special annotation
func (s *SerfPublisher) Publish(service v1.Service) (v1.Service, error) {
	newService := service.DeepCopy()
	cmd := exec.Command("/usr/sbin/avahi-ps", "publish", "TEST", "TEST", "9003", "TEST")
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Infof("cmd.Run() failed with %s\n", err)
	}
	s.logger.Infof("command \n%s\n", out)
	s.logger.Infof("HEMOS SIDO ENGAÑADOSSSS " + newService.ObjectMeta.Name)
	// newService, err := s.client.CoreV1().Services(newService.Namespace).Update(newService)

	return *newService, nil
}
