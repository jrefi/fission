/*
Copyright 2016 The Fission Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package buildermgr

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/fission/fission/pkg/crd"
	fetcherConfig "github.com/fission/fission/pkg/fetcher/config"
)

// Start the buildermgr service.
func Start(logger *zap.Logger, storageSvcUrl string, envBuilderNamespace string) error {
	bmLogger := logger.Named("builder_manager")

	fissionClient, kubernetesClient, _, err := crd.MakeFissionClient()
	if err != nil {
		return errors.Wrap(err, "failed to get fission or kubernetes client")
	}

	err = fissionClient.WaitForCRDs()
	if err != nil {
		return errors.Wrap(err, "error waiting for CRDs")
	}

	fetcherConfig, err := fetcherConfig.MakeFetcherConfig("/packages")
	if err != nil {
		return errors.Wrap(err, "error making fetcher config")
	}

	envWatcher := makeEnvironmentWatcher(bmLogger, fissionClient, kubernetesClient, fetcherConfig, envBuilderNamespace)
	go envWatcher.watchEnvironments()

	pkgWatcher := makePackageWatcher(bmLogger, fissionClient,
		kubernetesClient, envBuilderNamespace, storageSvcUrl)
	go pkgWatcher.watchPackages(fissionClient, kubernetesClient, envBuilderNamespace)

	select {}
}
