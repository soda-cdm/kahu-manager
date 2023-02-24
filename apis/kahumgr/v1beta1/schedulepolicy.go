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
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DaysType = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Sun":       time.Sunday,
	"Monday":    time.Monday,
	"Mon":       time.Monday,
	"Tuesday":   time.Tuesday,
	"Tue":       time.Tuesday,
	"Wednesday": time.Wednesday,
	"Wed":       time.Wednesday,
	"Thursday":  time.Thursday,
	"Thurs":     time.Thursday,
	"Friday":    time.Friday,
	"Fri":       time.Friday,
	"Saturday":  time.Saturday,
	"Sat":       time.Saturday,
}

const (
	layoutTime               = "15:09"
	HourlyPolicyType  string = "Hourly"
	DailyPolicyType   string = "Daily"
	WeeklyPolicyType  string = "Weekly"
	MonthlyPolicyType string = "Monthly"
)

func checkTimeFormat(policyTime string) error {
	_, err := time.Parse(layoutTime, policyTime)
	if err != nil {
		return fmt.Errorf("policyTime is: %s, err format :%v, you should provide the time"+
			" in the 0:00-23-59 format", policyTime, err)
	}
	return nil
}

// for every number of minutes after every hour the schedule will be triggered
// the cron example  25 * * * * (so every hour after 25 minutes triggers)
type HourlyPolicy struct {
	// Minutes when the policy should be triggered
	// +kubebuilder:validation:Maximum=59
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Minutes int `json:"minutes"`
	// +kubebuilder:validation:Maximum=256
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=24
	// +kubebuilder:validation:Optional
	MaxCopies int `json:"maxCopies"`
}

// Daily Policy contains the time in the day when the action should be triggered
// the cron example  20 16 * * *  (so every day at 16:20 Hrs triggers)
type DailyPolicy struct {
	// Time when the policy should be triggered
	// time eg 12:15
	// +kubebuilder:validation:Required
	Time string `json:"time"`
	// +kubebuilder:validation:Maximum=256
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=15
	// +kubebuilder:validation:Optional
	MaxCopies int `json:"maxCopies"`
}

func (d *DailyPolicy) CheckTimeFormat() error {
	return checkTimeFormat(d.Time)
}

// validates the DailyPolicy
func (d *DailyPolicy) Validate() error {
	return d.CheckTimeFormat()
}

// Weekly Policy contains the days and time  in a week when the action should be triggered
// the cron example  25 11 ? * (1,2) (so on Mon,Tues at 11:25 Hrs triggers)
type WeeklyPolicy struct {
	// Days of the week when the policy should be triggered.
	// Expected format is  specified in  DaysType as above
	// +kubebuilder:validation:Required
	Days []string `json:"days"`
	// +kubebuilder:validation:Required
	Time string `json:"time"`
	// +kubebuilder:validation:Maximum=256
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=4
	// +kubebuilder:validation:Optional
	MaxCopies int `json:"maxCopies"`
}

func (w *WeeklyPolicy) CheckTimeFormat() error {
	return checkTimeFormat(w.Time)
}

// validates the WeeklyPolicy
func (w *WeeklyPolicy) Validate() error {
	err := w.CheckTimeFormat()
	if err != nil {
		return err
	}
	var record1 []string
	var record2 []string
	strMap := make(map[time.Weekday]string)

	for _, day := range w.Days {
		// invalid format checking
		if _, exist1 := DaysType[day]; !exist1 {
			record1 = append(record1, day)
		}
		// duplicate checking
		if _, exists := strMap[DaysType[day]]; exists {
			record2 = append(record2, day)
		} else {
			strMap[DaysType[day]] = day
		}
	}

	if len(record1) > 0 {
		return fmt.Errorf("invalid day of the week: %v, in WeeklyPolicy you should provide the day"+
			" as Sunday, Sun, Monday, Mon, TuesDay, Tue, Wednesday, Wed, Thursday, Thurs, Friday, Fri,"+
			"Saturday, Sat", record1)
	}
	if len(record2) > 0 {
		return fmt.Errorf("Duplictae day of the week: %v, in WeeklyPolicy ", record2)
	}
	return nil
}

// Monthly Policy contains the dates and time  in a month when the action should be triggered
// the cron example  25 11 (1,5,8,11,18) * ?
// (so on given dates every month at 11:25 Hrs triggers)
type MonthlyPolicy struct {
	// Dates of the month when action should be triggered. If given date does not exist in a month then rollover to the next
	// date of month. Example 31 is specified then in Feb it will trigger on either 1st or 2nd March based on leap year or not.
	// +kubebuilder:validation:Required
	Dates []int `json:"dates"`
	// eg 12:15
	// +kubebuilder:validation:Required
	Time string `json:"time"`
	// +kubebuilder:validation:Maximum=256
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=12
	// +kubebuilder:validation:Optional
	MaxCopies int `json:"maxCopies"`
}

func (m *MonthlyPolicy) CheckTimeFormat() error {
	return checkTimeFormat(m.Time)
}

// validates the MonthlyPolicy
func (m *MonthlyPolicy) Validate() error {
	err := m.CheckTimeFormat()
	if err != nil {
		return err
	}
	var record []int
	for _, date := range m.Dates {
		// invalid format checking
		if date <= 0 || date > 31 {
			record = append(record, date)
		}
	}

	if len(record) > 0 {
		return fmt.Errorf("invalid date of the Month: %v, in MonthlyPolicy you should provide the date"+
			" as (0,31] only", record)
	}
	intMap := make(map[int]int)
	// validate for duplicate values
	for _, v := range m.Dates {
		// duplicate checking
		if _, exists := intMap[v]; exists {
			record = append(record, v)
		} else {
			intMap[v] = v
		}
	}

	if len(record) > 0 {
		return fmt.Errorf("duplicate date of the Month: %v, in MonthlyPolicy you should provide the date"+
			" as (0,31] only", record)
	}
	return nil
}

// SchedulePolicyspec
type SchedulePolicySpec struct {
	Hourly  *HourlyPolicy  `json:"hourly,omitempty"`
	Daily   *DailyPolicy   `json:"daily,omitempty"`
	Weekly  *WeeklyPolicy  `json:"weekly,omitempty"`
	Monthly *MonthlyPolicy `json:"monthly,omitempty"`
}

// SchedulePolicy is the Schema for the policy API
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
type SchedulePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SchedulePolicySpec `json:"spec,omitempty"`
}

// SchedulePolicyList contains a List of SchedulePolicy
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SchedulePolicyList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Items             []SchedulePolicy `json:"items"`
}
