package provider

type CronJobs struct {
	CronJobs []CronJobData `json:"cronjobs"`
	Metadata *Metadata     `json:"metadata"`
}

type CronJobData struct {
	Metadata CronJobMetadata `json:"metadata"`
	Spec     CronJobSpec     `json:"spec"`
	Status   CronJobStatus   `json:"status"`
}

type CronJobMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
}

type CronJobSpec struct {
	Schedule                   string          `json:"schedule"`
	ConcurrencyPolicy          string          `json:"concurrencyPolicy"`
	StartingDeadlineSeconds    *int64          `json:"startingDeadlineSeconds"`
	Suspend                    *bool           `json:"suspend"`
	JobTemplate                JobTemplateSpec `json:"jobTemplate"`
	SuccessfulJobsHistoryLimit *int32          `json:"successfulJobsHistoryLimit"`
	FailedJobsHistoryLimit     *int32          `json:"failedJobsHistoryLimit"`
}

type JobTemplateSpec struct {
	Metadata map[string]string `json:"metadata"`
	Spec     JobSpec           `json:"spec"`
}

type JobSpec struct {
	Template PodTemplateSpec `json:"template"`
}

type PodTemplateSpec struct {
	Metadata map[string]any `json:"metadata"`
	Spec     PodSpec        `json:"spec"`
}

type CronJobStatus struct {
	Active           []ObjectReference `json:"active"`
	LastScheduleTime string            `json:"lastScheduleTime"`
}

type ObjectReference struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
}
