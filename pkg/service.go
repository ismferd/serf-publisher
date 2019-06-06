package pkg

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os/exec"
	"strconv"
)

// SerfPublisherInterface for the SerfPublisher
type SerfPublisherInterface interface {
	Publish(service v1.Service) (v1.Service, error)
	Unpublish(name string) error
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

// Publish will add a new service trhought serf
func (s *SerfPublisher) Publish(service v1.Service) (v1.Service, error) {
	newService := service.DeepCopy()
	cmd := exec.Command("/usr/sbin/avahi-ps", "publish", newService.ObjectMeta.Name, "kubernetes", strconv.Itoa(int(newService.Spec.Ports[0].NodePort)), string(newService.Spec.Type))
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Infof("cmd.Run() failed with %s\n", err)
	}
	s.logger.Infof("command \n%s\n", out)
	return *newService, nil
}

// UnPublish will add a new service trhought serf
func (s *SerfPublisher) Unpublish(name string) error {
	options := metav1.GetOptions{}
	newService, _ := s.client.CoreV1().Services("cloudy").Get(name, options)

	cmd := exec.Command("/usr/sbin/avahi-ps", "unpublish", "kubernetes", strconv.Itoa(int(newService.Spec.Ports[0].NodePort)))
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Infof("cmd.Run() failed with %s\n", err)
	}
	s.logger.Infof("command \n%s\n", out)
	return err
}
