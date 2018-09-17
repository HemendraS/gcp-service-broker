// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// This file contains versioned models for the database so we
// can do proper tracking through gorm.
//
// If you need to change a model you MUST make a copy here and update the
// reference to the new model in db.go and add a migration path in the
// db_service package.

// ServiceBindingCredentialsV1 holds credentials returned to the users after
// binding to a service.
type ServiceBindingCredentialsV1 struct {
	gorm.Model

	OtherDetails string `gorm:"type:text"`

	ServiceId         string
	ServiceInstanceId string
	BindingId         string
}

// TableName returns a consistent table name (`service_binding_credentials`) for
// gorm so multiple structs from different versions of the database all operate
// on the same table.
func (ServiceBindingCredentialsV1) TableName() string {
	return "service_binding_credentials"
}

// ServiceInstanceDetailsV1 holds information about provisioned services.
type ServiceInstanceDetailsV1 struct {
	ID        string `gorm:"primary_key;type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name         string
	Location     string
	Url          string
	OtherDetails string `gorm:"type:text"`

	ServiceId        string
	PlanId           string
	SpaceGuid        string
	OrganizationGuid string
}

// TableName returns a consistent table name (`service_instance_details`) for
// gorm so multiple structs from different versions of the database all operate
// on the same table.
func (ServiceInstanceDetailsV1) TableName() string {
	return "service_instance_details"
}

// ProvisionRequestDetailsV1 holds user-defined properties passed to a call
// to provision a service.
type ProvisionRequestDetailsV1 struct {
	gorm.Model

	ServiceInstanceId string
	// is a json.Marshal of models.ProvisionDetails
	RequestDetails string
}

// TableName returns a consistent table name (`provision_request_details`) for
// gorm so multiple structs from different versions of the database all operate
// on the same table.
func (ProvisionRequestDetailsV1) TableName() string {
	return "provision_request_details"
}

// MigrationV1 represents the mgirations table. It holds a monotonically
// increasing number that gets incremented with every database schema revision.
type MigrationV1 struct {
	gorm.Model

	MigrationId int `gorm:"type:int(10)"`
}

// TableName returns a consistent table name (`migrations`) for gorm so
// multiple structs from different versions of the database all operate on the
// same table.
func (MigrationV1) TableName() string {
	return "migrations"
}

// CloudOperationV1 holds information about the status of Google Cloud
// long-running operations.
type CloudOperationV1 struct {
	gorm.Model

	Name          string
	Status        string
	OperationType string
	ErrorMessage  string `gorm:"type:text"`
	InsertTime    string
	StartTime     string
	TargetId      string
	TargetLink    string

	ServiceId         string
	ServiceInstanceId string
}

// TableName returns a consistent table name (`cloud_operations`) for gorm so
// multiple structs from different versions of the database all operate on the
// same table.
func (CloudOperationV1) TableName() string {
	return "cloud_operations"
}

// PlanDetailsV1 is a table that was deprecated in favor of using Environment
// variables. It only remains for ORM migrations and the ability for existing
// users to export their plans.
type PlanDetailsV1 struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	ServiceId string
	Name      string
	Features  string `sql:"type:text"`
}

// TableName returns a consistent table name (`plan_details`) for gorm so
// multiple structs from different versions of the database all operate on the
// same table.
func (PlanDetailsV1) TableName() string {
	return "plan_details"
}