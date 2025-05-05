package provider

type Pods struct {
	Pods     []PodData `json:"pods"`
	Metadata *Metadata `json:"metadata"`
}

type PodData struct {
	Metadata PodMetadata `json:"metadata"`
	Spec     PodSpec     `json:"spec"`
	Status   PodStatus   `json:"status"`
}

type PodMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	OwnerReferences   []struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	} `json:"ownerReferences"`
}

type PodSpec struct {
	Containers         []Container     `json:"containers"`
	NodeName           string          `json:"nodeName"`
	RestartPolicy      string          `json:"restartPolicy"`
	ServiceAccountName string          `json:"serviceAccountName"`
	HostIP             string          `json:"hostIP"`
	SecurityContext    SecurityContext `json:"securityContext"`
}

type PodStatus struct {
	Phase             string            `json:"phase"`
	Conditions        []PodCondition    `json:"conditions"`
	HostIP            string            `json:"hostIP"`
	PodIP             string            `json:"podIP"`
	StartTime         string            `json:"startTime"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses"`
}

type PodCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastProbeTime      string `json:"lastProbeTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

type ContainerStatus struct {
	Name         string `json:"name"`
	State        State  `json:"state"`
	LastState    State  `json:"lastState"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartCount"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	ContainerID  string `json:"containerID"`
	Started      bool   `json:"started"`
}

type State struct {
	Waiting    *StateDetails `json:"waiting,omitempty"`
	Running    *StateDetails `json:"running,omitempty"`
	Terminated *StateDetails `json:"terminated,omitempty"`
}

type StateDetails struct {
	Reason   string `json:"reason"`
	Message  string `json:"message"`
	ExitCode int    `json:"exitCode,omitempty"`
	Signal   int    `json:"signal,omitempty"`
}
