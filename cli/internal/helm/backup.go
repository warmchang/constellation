package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/edgelesssys/constellation/v2/internal/constants"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

// TODO: kubectl get crds -n kube-system -o yaml > ${workspace}/constellation-upgrade/backups/crds/...
// TODO: for name in ^ { kubectl get $name -n kube-system -o yaml > ${workspace}/constellation-upgrade/backups/...}
const (
	crdBackupFolder = "constellation-upgrade/backups/crds/"
	backupFolder    = "constellation-upgrade/backups/"
)

func (c *Client) backupCRDs(ctx context.Context) ([]apiextensionsv1.CustomResourceDefinition, error) {
	crds, err := c.kubectl.GetCRDs(ctx, constants.HelmNamespace)
	if err != nil {
		return nil, fmt.Errorf("getting CRDs: %w", err)
	}

	if err := os.MkdirAll(crdBackupFolder, 0o700); err != nil {
		return nil, fmt.Errorf("creating backup dir: %w", err)
	}
	for i := range crds {
		path := filepath.Join(crdBackupFolder, crds[i].Name+".yaml")

		// We have to manually set kind/apiversion because of a long-standing limitation of the API:
		// https://github.com/kubernetes/kubernetes/issues/3030#issuecomment-67543738
		// The comment states that kind/version are encoded in the type.
		// The package holding the CRD type encodes the version.
		crds[i].Kind = "CustomResourceDefinition"
		crds[i].APIVersion = "apiextensions.k8s.io/v1"

		yamlBytes, err := yaml.Marshal(crds[i])
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, yamlBytes, 0o644); err != nil {
			return nil, err
		}

		c.log.Debugf("created backup crd: %s", path)
	}
	return crds, nil
}

func (c *Client) backupCRs(crds []apiextensionsv1.CustomResourceDefinition) error {
	for _, crd := range crds {
		tmp := crd.Name
		c.kubectl.GetCRs(ctx, crd.Name)
		c.log.Debugf("name: %s", tmp)
	}
	return nil
}
