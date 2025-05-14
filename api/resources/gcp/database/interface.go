package sql

import (
	googleapi "google.golang.org/api/googleapi"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type SQLAdmin interface {
	List(project string) Caller
}

type Caller interface {
	Do(opts ...googleapi.CallOption) (*sqladmin.InstancesListResponse, error)
}
