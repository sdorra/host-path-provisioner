package storage

import (
	"log"
	"os"
	"path"

	"github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pkg/errors"
)

const (
	NAME = "sdorra.org/host-path"
)

func NewHostPathProvisioner(directory string) *HostPathProvisioner {
	return &HostPathProvisioner{
		directory: directory,
	}
}

type HostPathProvisioner struct {
	directory string
}

func (p *HostPathProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	log.Printf("provisioning volume %v", options)

	volumePath := path.Join(p.directory, options.PVName)
	log.Printf("create volume %s at %v", options.PVName, volumePath)
	if err := os.MkdirAll(volumePath, 0777); err != nil {
		return nil, err
	}

	if err := os.Chmod(volumePath, 0777); err != nil {
		return nil, err
	}

	// fix subpath container creation
  // https://github.com/kubernetes/kubernetes/issues/66583
	hostPathType := v1.HostPathDirectoryOrCreate

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: volumePath,
					Type: &hostPathType,
				},
			},
		},
	}
	return pv, nil
}

func (p *HostPathProvisioner) Delete(volume *v1.PersistentVolume) error {
	volumePath := path.Join(p.directory, volume.Name)
	if err := os.RemoveAll(volumePath); err != nil {
		return errors.Wrapf(err, "removing host path volume %s", volumePath)
	}

	return nil
}
