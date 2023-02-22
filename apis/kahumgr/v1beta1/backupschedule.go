/*
Copyright 2023 The SODA Authors.

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

package v1beta1

import (
	kahubk "github.com/soda-cdm/kahu/apis/kahu/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConcurrencyPolicy describes how the BackupSchedule will be handled.
// Only one of the following concurrent policies may be specified.
// If none of the following policies is specified, the default one
// is ForbidConcurrent.
type ConcurrencyPolicy string

const (
	// AllowConcurrent allows Backup to run concurrently.
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent forbids concurrent runs, skipping next run if previous
	// hasn't finished yet.
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent cancels currently running Backup and replaces it with a new one.
	// ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

// ReclaimPolicy tells about reclamation of the backup. It can be either delete or retain
type BackupScheduleReclaimPolicyType struct {
	// +optional
	ReclaimPolicyDelete string `json:"reclaimPolicyDelete,omitempty"`

	// +optional
	ReclaimPolicyRetain string `json:"reclaimPolicyRetain,omitempty"`
}

type BackupScheduleSpec struct {
	// optional, name of the SchedulePolicy CR
	// if empty considered as manual trigger otherwise scheduled based backup will be taken
	BackupPolicyName string `json:"backupPolicyName"`
	// ReclaimPolicy tells about reclamation of the backup. It can be either delete or retain
	// +kubebuilder:default= retain
	// +kubebuilder:validation:Optional
	ReclaimPolicy BackupScheduleReclaimPolicyType `json:"reclaimPolicy,omitempty"`
	// Enable tells whether  Scheduled Backup should be started or stopped
	// +optional
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	Enable bool `json:"enable,omitempty"`
	// +kubebuilder:validation:Maximum=5
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=3
	// +kubebuilder:validation:Optional
	MaxRetriesOnFailure int `json:"maxRetriesOnFailure"`
	// Optional deadline in seconds for starting the Backup if it  misses
	// scheduled time for any reason.
	// +optional
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`
	// Specifies how to treat concurrent executions of a Backup.
	// Valid values are:
	// - "Allow": allows Backups to run concurrently;
	// - "Forbid"(default): forbids concurrent runs, skipping next run if previous run hasn't finished yet.
	// +optional
	ConcurrentPolicy ConcurrencyPolicy `json:"concurrentPolicy,omitempty"`
	// this Backup spec
	BackupTemplate kahubk.BackupSpec `json:"template,omitempty"`
}

type ScheduleStatus string
type ExecutionStatus string

const (
	SchedulePending  ScheduleStatus = "Pending"
	ScheduleActive   ScheduleStatus = "Active"
	ScheduleInActive ScheduleStatus = "InActive"
	ScheduleFailed   ScheduleStatus = "Failed"
	ScheduleDeleting ScheduleStatus = "Deleting"

	ExecutionSuccess    ExecutionStatus = "Success"
	ExecutionInProgress ExecutionStatus = "InProgress"
	ExecutionFailure    ExecutionStatus = "Failed"
)

type StatusInfo struct {
	BackupName          string          `json:"backupName"`
	ExecStatus          ExecutionStatus `json:"execStatus"`
	StartTimestamp      metav1.Time     `json:"startTimestamp"`
	CompletionTimestamp metav1.Time     `json:"completionTimestamp"`
}

type BackupScheduleStatus struct {
	// latest 10 scheduled backup status is stored
	RecentStatusInfo    []StatusInfo    `json:"recentStatusInfo"`
	LastBackupName      string          `json:"lastBackupName"`
	LastExecutionStatus ExecutionStatus `json:"lastExecutionStatus"`
	// LastStartTimestamp is defines time when Schedule created the backup
	LastStartTimestamp metav1.Time `json:"lastStartTimestamp"`
	// LastCompletionTimestamp is defines time when backup completed
	LastCompletionTimestamp metav1.Time    `json:"lastCompletionTimestamp"`
	SchedStatus             ScheduleStatus `json:"schedStatus"`
	// the created backup crd status used to identify the completed or not
	BackupStatus kahubk.BackupState `json:"backupStatus"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="LastBackupName",type=string,JSONPath=`.status.lastBackupName`,description="Name of the recent backup triggered based on this backupschedule."
// +kubebuilder:printcolumn:name="BackupPolicyName",type=string,JSONPath=`.spec.backupPolicyName`,description=" schedule policy Name."
// +kubebuilder:printcolumn:name="Enable",type=boolean,JSONPath=`.spec.enable`,description="Indicates the backup trigger is enabled or disabled."
type BackupSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BackupScheduleSpec   `json:"spec,omitempty"`
	Status            BackupScheduleStatus `json:"status,omitempty"`
}

// BackupScheduleList contains a List of BackupSchedule
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BackupScheduleList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []BackupSchedule `json:"items"`
}
