package pkg

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Handler will receive objects whenever they get added or deleted from the k8s API.
type Handler struct {
	serfPublisher SerfPublisherInterface
}

// Add is called when a k8s object is created.
func (h *Handler) Add(obj runtime.Object) error {
	_, err := h.serfPublisher.Publish(*obj.(*v1.Service))
	return err
}

// Delete is called when a k8s object is deleted.
func (h *Handler) Delete(s string) error {
	return nil
}

// NewHandler returns a new Handler to handle Deployments created/updated/deleted.
func NewHandler(serfPublisher SerfPublisherInterface) *Handler {
	return &Handler{
		serfPublisher: serfPublisher,
	}
}
