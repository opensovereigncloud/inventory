// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package crd

import (
	"context"

	metalv1alpha4 "github.com/ironcore-dev/metal/api/v1alpha1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	CSonicNamespace = "onmetal.de"
)

type KubeAPISaverSvc struct {
	client client.Client
}

func NewKubeAPISaverSvc(kubeconfig string, namespace string) (SaverSvc, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read kubeconfig from path %s", kubeconfig)
	}

	if err := metalv1alpha4.AddToScheme(scheme.Scheme); err != nil {
		return nil, errors.Wrap(err, "unable to add registered types to client scheme")
	}

	// clientset, err := clientv1alpha1.NewForConfig(config)
	cl, err := client.New(config, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, errors.Wrap(err, "unable to build clientset from config")
	}

	// client := clientset.Inventories(namespace)

	return &KubeAPISaverSvc{
		client: cl,
	}, nil
}

func (s *KubeAPISaverSvc) Save(inv *metalv1alpha4.Inventory) error {
	err := s.client.Create(context.Background(), inv)
	if err == nil {
		return nil
	}
	if !apierrors.IsAlreadyExists(err) {
		return errors.Wrap(err, "unhandled error on creation")
	}

	existing := &metalv1alpha4.Inventory{}
	err = s.client.Get(context.Background(), types.NamespacedName{
		Namespace: "",
		Name:      inv.Name,
	}, existing)
	if err != nil {
		return errors.Wrap(err, "unable to get resource")
	}

	existing.Spec = inv.Spec

	if err = s.client.Update(context.Background(), existing); err != nil {
		return errors.Wrap(err, "unhandled error on update")
	}

	return nil
}

func (s *KubeAPISaverSvc) Patch(_ string, _ interface{}) error {
	// patchData, err := json.Marshal(patch)
	// if err != nil {
	// 	return errors.Wrap(err, "unable to marshal inventory")
	// }
	// fmt.Println(string(patchData))
	// err = s.client.Patch(context.Background(), name, types.MergePatchType, patchData, metav1.PatchOptions{})
	// if err != nil {
	// 	return errors.Wrap(err, "unable to patch inventory")
	// }
	return nil
}
