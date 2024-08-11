package patch

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type RequestOpBody struct {
	EndpointPath string   `json:"endpoint"`
	Objects      []Object `json:"objects"`
}

type Object struct {
	Op                string      `json:"op"`
	RequestObjectType string      `json:"objectType"`
	TargetObjectType  string      `json:"targetObjectType"`
	RequestSpec       interface{} `json:"spec"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

// Constant for the target request	object type
const (
	CONTAINER = "container"
	POD       = "pod"
	VOLUME    = "volume"
)

// Constant for the target object type
const (
	DEPLOYMENT  = "deployment"
	STATEFULSET = "statefulset"
)
