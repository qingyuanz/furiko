//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
 * Copyright 2022 The Furiko Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BoolOptionConfig) DeepCopyInto(out *BoolOptionConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BoolOptionConfig.
func (in *BoolOptionConfig) DeepCopy() *BoolOptionConfig {
	if in == nil {
		return nil
	}
	out := new(BoolOptionConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConcurrencySpec) DeepCopyInto(out *ConcurrencySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConcurrencySpec.
func (in *ConcurrencySpec) DeepCopy() *ConcurrencySpec {
	if in == nil {
		return nil
	}
	out := new(ConcurrencySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronSchedule) DeepCopyInto(out *CronSchedule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronSchedule.
func (in *CronSchedule) DeepCopy() *CronSchedule {
	if in == nil {
		return nil
	}
	out := new(CronSchedule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DateOptionConfig) DeepCopyInto(out *DateOptionConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DateOptionConfig.
func (in *DateOptionConfig) DeepCopy() *DateOptionConfig {
	if in == nil {
		return nil
	}
	out := new(DateOptionConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Job) DeepCopyInto(out *Job) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Job.
func (in *Job) DeepCopy() *Job {
	if in == nil {
		return nil
	}
	out := new(Job)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Job) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobCondition) DeepCopyInto(out *JobCondition) {
	*out = *in
	if in.Queueing != nil {
		in, out := &in.Queueing, &out.Queueing
		*out = new(JobConditionQueueing)
		**out = **in
	}
	if in.Waiting != nil {
		in, out := &in.Waiting, &out.Waiting
		*out = new(JobConditionWaiting)
		**out = **in
	}
	if in.Running != nil {
		in, out := &in.Running, &out.Running
		*out = new(JobConditionRunning)
		(*in).DeepCopyInto(*out)
	}
	if in.Finished != nil {
		in, out := &in.Finished, &out.Finished
		*out = new(JobConditionFinished)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobCondition.
func (in *JobCondition) DeepCopy() *JobCondition {
	if in == nil {
		return nil
	}
	out := new(JobCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConditionFinished) DeepCopyInto(out *JobConditionFinished) {
	*out = *in
	if in.LatestCreationTimestamp != nil {
		in, out := &in.LatestCreationTimestamp, &out.LatestCreationTimestamp
		*out = (*in).DeepCopy()
	}
	if in.LatestRunningTimestamp != nil {
		in, out := &in.LatestRunningTimestamp, &out.LatestRunningTimestamp
		*out = (*in).DeepCopy()
	}
	in.FinishTimestamp.DeepCopyInto(&out.FinishTimestamp)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConditionFinished.
func (in *JobConditionFinished) DeepCopy() *JobConditionFinished {
	if in == nil {
		return nil
	}
	out := new(JobConditionFinished)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConditionQueueing) DeepCopyInto(out *JobConditionQueueing) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConditionQueueing.
func (in *JobConditionQueueing) DeepCopy() *JobConditionQueueing {
	if in == nil {
		return nil
	}
	out := new(JobConditionQueueing)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConditionRunning) DeepCopyInto(out *JobConditionRunning) {
	*out = *in
	in.LatestCreationTimestamp.DeepCopyInto(&out.LatestCreationTimestamp)
	in.LatestRunningTimestamp.DeepCopyInto(&out.LatestRunningTimestamp)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConditionRunning.
func (in *JobConditionRunning) DeepCopy() *JobConditionRunning {
	if in == nil {
		return nil
	}
	out := new(JobConditionRunning)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConditionWaiting) DeepCopyInto(out *JobConditionWaiting) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConditionWaiting.
func (in *JobConditionWaiting) DeepCopy() *JobConditionWaiting {
	if in == nil {
		return nil
	}
	out := new(JobConditionWaiting)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConfig) DeepCopyInto(out *JobConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConfig.
func (in *JobConfig) DeepCopy() *JobConfig {
	if in == nil {
		return nil
	}
	out := new(JobConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JobConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConfigList) DeepCopyInto(out *JobConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]JobConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConfigList.
func (in *JobConfigList) DeepCopy() *JobConfigList {
	if in == nil {
		return nil
	}
	out := new(JobConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JobConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConfigSpec) DeepCopyInto(out *JobConfigSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
	out.Concurrency = in.Concurrency
	if in.Schedule != nil {
		in, out := &in.Schedule, &out.Schedule
		*out = new(ScheduleSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Option != nil {
		in, out := &in.Option, &out.Option
		*out = new(OptionSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConfigSpec.
func (in *JobConfigSpec) DeepCopy() *JobConfigSpec {
	if in == nil {
		return nil
	}
	out := new(JobConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobConfigStatus) DeepCopyInto(out *JobConfigStatus) {
	*out = *in
	if in.QueuedJobs != nil {
		in, out := &in.QueuedJobs, &out.QueuedJobs
		*out = make([]JobReference, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ActiveJobs != nil {
		in, out := &in.ActiveJobs, &out.ActiveJobs
		*out = make([]JobReference, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LastScheduled != nil {
		in, out := &in.LastScheduled, &out.LastScheduled
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobConfigStatus.
func (in *JobConfigStatus) DeepCopy() *JobConfigStatus {
	if in == nil {
		return nil
	}
	out := new(JobConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobList) DeepCopyInto(out *JobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Job, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobList.
func (in *JobList) DeepCopy() *JobList {
	if in == nil {
		return nil
	}
	out := new(JobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobReference) DeepCopyInto(out *JobReference) {
	*out = *in
	in.CreationTimestamp.DeepCopyInto(&out.CreationTimestamp)
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobReference.
func (in *JobReference) DeepCopy() *JobReference {
	if in == nil {
		return nil
	}
	out := new(JobReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobSpec) DeepCopyInto(out *JobSpec) {
	*out = *in
	if in.StartPolicy != nil {
		in, out := &in.StartPolicy, &out.StartPolicy
		*out = new(StartPolicySpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(JobTemplate)
		(*in).DeepCopyInto(*out)
	}
	if in.Substitutions != nil {
		in, out := &in.Substitutions, &out.Substitutions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.KillTimestamp != nil {
		in, out := &in.KillTimestamp, &out.KillTimestamp
		*out = (*in).DeepCopy()
	}
	if in.TTLSecondsAfterFinished != nil {
		in, out := &in.TTLSecondsAfterFinished, &out.TTLSecondsAfterFinished
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobSpec.
func (in *JobSpec) DeepCopy() *JobSpec {
	if in == nil {
		return nil
	}
	out := new(JobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobStatus) DeepCopyInto(out *JobStatus) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
	if in.Tasks != nil {
		in, out := &in.Tasks, &out.Tasks
		*out = make([]TaskRef, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ParallelStatus != nil {
		in, out := &in.ParallelStatus, &out.ParallelStatus
		*out = new(ParallelStatus)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobStatus.
func (in *JobStatus) DeepCopy() *JobStatus {
	if in == nil {
		return nil
	}
	out := new(JobStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobTemplate) DeepCopyInto(out *JobTemplate) {
	*out = *in
	in.TaskTemplate.DeepCopyInto(&out.TaskTemplate)
	if in.Parallelism != nil {
		in, out := &in.Parallelism, &out.Parallelism
		*out = new(ParallelismSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.MaxAttempts != nil {
		in, out := &in.MaxAttempts, &out.MaxAttempts
		*out = new(int64)
		**out = **in
	}
	if in.RetryDelaySeconds != nil {
		in, out := &in.RetryDelaySeconds, &out.RetryDelaySeconds
		*out = new(int64)
		**out = **in
	}
	if in.TaskPendingTimeoutSeconds != nil {
		in, out := &in.TaskPendingTimeoutSeconds, &out.TaskPendingTimeoutSeconds
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobTemplate.
func (in *JobTemplate) DeepCopy() *JobTemplate {
	if in == nil {
		return nil
	}
	out := new(JobTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobTemplateSpec) DeepCopyInto(out *JobTemplateSpec) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobTemplateSpec.
func (in *JobTemplateSpec) DeepCopy() *JobTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(JobTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiOptionConfig) DeepCopyInto(out *MultiOptionConfig) {
	*out = *in
	if in.Default != nil {
		in, out := &in.Default, &out.Default
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultiOptionConfig.
func (in *MultiOptionConfig) DeepCopy() *MultiOptionConfig {
	if in == nil {
		return nil
	}
	out := new(MultiOptionConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Option) DeepCopyInto(out *Option) {
	*out = *in
	if in.Bool != nil {
		in, out := &in.Bool, &out.Bool
		*out = new(BoolOptionConfig)
		**out = **in
	}
	if in.String != nil {
		in, out := &in.String, &out.String
		*out = new(StringOptionConfig)
		**out = **in
	}
	if in.Select != nil {
		in, out := &in.Select, &out.Select
		*out = new(SelectOptionConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Multi != nil {
		in, out := &in.Multi, &out.Multi
		*out = new(MultiOptionConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Date != nil {
		in, out := &in.Date, &out.Date
		*out = new(DateOptionConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Option.
func (in *Option) DeepCopy() *Option {
	if in == nil {
		return nil
	}
	out := new(Option)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OptionSpec) DeepCopyInto(out *OptionSpec) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make([]Option, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OptionSpec.
func (in *OptionSpec) DeepCopy() *OptionSpec {
	if in == nil {
		return nil
	}
	out := new(OptionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParallelIndex) DeepCopyInto(out *ParallelIndex) {
	*out = *in
	if in.IndexNumber != nil {
		in, out := &in.IndexNumber, &out.IndexNumber
		*out = new(int64)
		**out = **in
	}
	if in.MatrixValues != nil {
		in, out := &in.MatrixValues, &out.MatrixValues
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParallelIndex.
func (in *ParallelIndex) DeepCopy() *ParallelIndex {
	if in == nil {
		return nil
	}
	out := new(ParallelIndex)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParallelStatus) DeepCopyInto(out *ParallelStatus) {
	*out = *in
	if in.Successful != nil {
		in, out := &in.Successful, &out.Successful
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParallelStatus.
func (in *ParallelStatus) DeepCopy() *ParallelStatus {
	if in == nil {
		return nil
	}
	out := new(ParallelStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ParallelismSpec) DeepCopyInto(out *ParallelismSpec) {
	*out = *in
	if in.WithCount != nil {
		in, out := &in.WithCount, &out.WithCount
		*out = new(int64)
		**out = **in
	}
	if in.WithKeys != nil {
		in, out := &in.WithKeys, &out.WithKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.WithMatrix != nil {
		in, out := &in.WithMatrix, &out.WithMatrix
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ParallelismSpec.
func (in *ParallelismSpec) DeepCopy() *ParallelismSpec {
	if in == nil {
		return nil
	}
	out := new(ParallelismSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodTemplateSpec) DeepCopyInto(out *PodTemplateSpec) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodTemplateSpec.
func (in *PodTemplateSpec) DeepCopy() *PodTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(PodTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleContraints) DeepCopyInto(out *ScheduleContraints) {
	*out = *in
	if in.NotBefore != nil {
		in, out := &in.NotBefore, &out.NotBefore
		*out = (*in).DeepCopy()
	}
	if in.NotAfter != nil {
		in, out := &in.NotAfter, &out.NotAfter
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleContraints.
func (in *ScheduleContraints) DeepCopy() *ScheduleContraints {
	if in == nil {
		return nil
	}
	out := new(ScheduleContraints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleSpec) DeepCopyInto(out *ScheduleSpec) {
	*out = *in
	if in.Cron != nil {
		in, out := &in.Cron, &out.Cron
		*out = new(CronSchedule)
		**out = **in
	}
	if in.Constraints != nil {
		in, out := &in.Constraints, &out.Constraints
		*out = new(ScheduleContraints)
		(*in).DeepCopyInto(*out)
	}
	if in.LastUpdated != nil {
		in, out := &in.LastUpdated, &out.LastUpdated
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleSpec.
func (in *ScheduleSpec) DeepCopy() *ScheduleSpec {
	if in == nil {
		return nil
	}
	out := new(ScheduleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SelectOptionConfig) DeepCopyInto(out *SelectOptionConfig) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SelectOptionConfig.
func (in *SelectOptionConfig) DeepCopy() *SelectOptionConfig {
	if in == nil {
		return nil
	}
	out := new(SelectOptionConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StartPolicySpec) DeepCopyInto(out *StartPolicySpec) {
	*out = *in
	if in.StartAfter != nil {
		in, out := &in.StartAfter, &out.StartAfter
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StartPolicySpec.
func (in *StartPolicySpec) DeepCopy() *StartPolicySpec {
	if in == nil {
		return nil
	}
	out := new(StartPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StringOptionConfig) DeepCopyInto(out *StringOptionConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StringOptionConfig.
func (in *StringOptionConfig) DeepCopy() *StringOptionConfig {
	if in == nil {
		return nil
	}
	out := new(StringOptionConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskContainerState) DeepCopyInto(out *TaskContainerState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskContainerState.
func (in *TaskContainerState) DeepCopy() *TaskContainerState {
	if in == nil {
		return nil
	}
	out := new(TaskContainerState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskRef) DeepCopyInto(out *TaskRef) {
	*out = *in
	in.CreationTimestamp.DeepCopyInto(&out.CreationTimestamp)
	if in.RunningTimestamp != nil {
		in, out := &in.RunningTimestamp, &out.RunningTimestamp
		*out = (*in).DeepCopy()
	}
	if in.FinishTimestamp != nil {
		in, out := &in.FinishTimestamp, &out.FinishTimestamp
		*out = (*in).DeepCopy()
	}
	if in.ParallelIndex != nil {
		in, out := &in.ParallelIndex, &out.ParallelIndex
		*out = new(ParallelIndex)
		(*in).DeepCopyInto(*out)
	}
	out.Status = in.Status
	if in.DeletedStatus != nil {
		in, out := &in.DeletedStatus, &out.DeletedStatus
		*out = new(TaskStatus)
		**out = **in
	}
	if in.ContainerStates != nil {
		in, out := &in.ContainerStates, &out.ContainerStates
		*out = make([]TaskContainerState, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskRef.
func (in *TaskRef) DeepCopy() *TaskRef {
	if in == nil {
		return nil
	}
	out := new(TaskRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskStatus) DeepCopyInto(out *TaskStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskStatus.
func (in *TaskStatus) DeepCopy() *TaskStatus {
	if in == nil {
		return nil
	}
	out := new(TaskStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TaskTemplate) DeepCopyInto(out *TaskTemplate) {
	*out = *in
	if in.Pod != nil {
		in, out := &in.Pod, &out.Pod
		*out = new(PodTemplateSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TaskTemplate.
func (in *TaskTemplate) DeepCopy() *TaskTemplate {
	if in == nil {
		return nil
	}
	out := new(TaskTemplate)
	in.DeepCopyInto(out)
	return out
}
