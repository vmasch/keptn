package helm

import (
	"helm.sh/helm/v3/pkg/chart"
)

// HelmExecutor is an interface for Helm operations
type HelmExecutor interface {
	UpgradeChart(ch *chart.Chart, releaseName, namespace string, vals map[string]interface{}) error
}
