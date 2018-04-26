package gube

import (
	"strings"

	"github.com/gardener/gardener/pkg/operation/common"
	corev1 "k8s.io/api/core/v1"
)

const gardenprefix = "garden-"

func GetProjectNameFromNamespaceManifest(m *corev1.Namespace) string {
	name, ok := m.GetLabels()[common.ProjectName]
	if !ok {
		name = m.GetName()
		if strings.HasPrefix(name, gardenprefix) {
			name = name[len(gardenprefix):]
		}
	}
	return name
}
