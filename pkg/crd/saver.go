// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package crd

import apiv1alpha1 "github.com/ironcore-dev/metal/api/v1alpha1"

type SaverSvc interface {
	Save(inv *apiv1alpha1.Inventory) error
}
