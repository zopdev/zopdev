package provider

type Deployments struct {
	Deployments []DeploymentData `json:"deployments"`
	Metadata    *Metadata        `json:"metadata"`
}

type DeploymentData struct {
	Metadata ItemMetadata `json:"metadata"`
	Spec     ItemSpec     `json:"spec"`
	Status   DepStatus    `json:"status"`
}

type ItemMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	Generation        int64             `json:"generation"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	ManagedFields     []ManagedField    `json:"managedFields"`
}

type ManagedField struct {
	Manager     string  `json:"manager"`
	Operation   string  `json:"operation"`
	APIVersion  string  `json:"apiVersion"`
	Time        string  `json:"time"`
	FieldsType  string  `json:"fieldsType"`
	Subresource *string `json:"subresource,omitempty"`
}

type ItemSpec struct {
	Replicas                int64    `json:"replicas"`
	Selector                Selector `json:"selector"`
	Template                Template `json:"template"`
	Strategy                Strategy `json:"strategy"`
	RevisionHistoryLimit    int64    `json:"revisionHistoryLimit"`
	ProgressDeadlineSeconds int64    `json:"progressDeadlineSeconds"`
}

type Selector struct {
	MatchLabels MatchLabelsClass `json:"matchLabels"`
}

type MatchLabelsClass struct {
	AppKubernetesIoInstance string `json:"app.kubernetes.io/instance"`
	AppKubernetesIoName     string `json:"app.kubernetes.io/name"`
}

type Strategy struct {
	Type          string        `json:"type"`
	RollingUpdate RollingUpdate `json:"rollingUpdate"`
}

type RollingUpdate struct {
	MaxUnavailable any `json:"maxUnavailable"`
	MaxSurge       any `json:"maxSurge"`
}

type Template struct {
	Metadata TemplateMetadata `json:"metadata"`
	Spec     TemplateSpec     `json:"spec"`
}

type TemplateMetadata struct {
	CreationTimestamp interface{}      `json:"creationTimestamp"`
	Labels            MatchLabelsClass `json:"labels"`
}

type TemplateSpec struct {
	Containers                    []Container     `json:"containers"`
	RestartPolicy                 string          `json:"restartPolicy"`
	TerminationGracePeriodSeconds int64           `json:"terminationGracePeriodSeconds"`
	DNSPolicy                     string          `json:"dnsPolicy"`
	ServiceAccountName            string          `json:"serviceAccountName"`
	ServiceAccount                string          `json:"serviceAccount"`
	SecurityContext               SecurityContext `json:"securityContext"`
	SchedulerName                 string          `json:"schedulerName"`
}

type Container struct {
	Name                     string     `json:"name"`
	Image                    string     `json:"image"`
	Ports                    []DepPorts `json:"ports"`
	Env                      []Env      `json:"env"`
	Resources                Resources  `json:"resources"`
	LivenessProbe            Probe      `json:"livenessProbe"`
	ReadinessProbe           Probe      `json:"readinessProbe"`
	StartupProbe             *Probe     `json:"startupProbe,omitempty"`
	TerminationMessagePath   string     `json:"terminationMessagePath"`
	TerminationMessagePolicy string     `json:"terminationMessagePolicy"`
	ImagePullPolicy          string     `json:"imagePullPolicy"`
	Status                   string     `json:"status"`
	Command                  any        `json:"command"`
	Args                     any        `json:"args"`
	VolumeMounts             any        `json:"volumeMounts"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Probe struct {
	HTTPGet             HTTPGet `json:"httpGet"`
	InitialDelaySeconds int64   `json:"initialDelaySeconds"`
	TimeoutSeconds      int64   `json:"timeoutSeconds"`
	PeriodSeconds       int64   `json:"periodSeconds"`
	SuccessThreshold    int64   `json:"successThreshold"`
	FailureThreshold    int64   `json:"failureThreshold"`
}

type HTTPGet struct {
	Path   string `json:"path"`
	Port   any    `json:"port"`
	Scheme string `json:"scheme"`
}

type DepPorts struct {
	Name          string `json:"name"`
	ContainerPort int64  `json:"containerPort"`
	Protocol      string `json:"protocol"`
}

type Resources struct {
	Limits   Limits `json:"limits"`
	Requests Limits `json:"requests"`
}

type Limits struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type SecurityContext struct {
}

type DepStatus struct {
	ObservedGeneration int64       `json:"observedGeneration"`
	Replicas           int64       `json:"replicas"`
	UpdatedReplicas    int64       `json:"updatedReplicas"`
	ReadyReplicas      int64       `json:"readyReplicas"`
	AvailableReplicas  int64       `json:"availableReplicas"`
	Conditions         []Condition `json:"conditions"`
}

type Condition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastUpdateTime     string `json:"lastUpdateTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
}
