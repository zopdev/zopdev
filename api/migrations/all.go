// This is auto-generated file using 'gofr migrate' tool. DO NOT EDIT.
package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{

		20241209162239: createCloudAccountTable(),
		20241211121223: createAPPlicationTable(),
		20241211121308: createEnvironmentTable(),
		20241211121841: createDeploymentSpaceTable(),
		20241212162207: createClusterTable(),
	}
}
