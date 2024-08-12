package patch

type PatchBase struct {
	EndpointPath string   `json:"endpoint"`
	Objects      []Object `json:"objects"`
}

type Object struct {
	Op                string      `json:"op"`
	RequestObjectType string      `json:"objectType"`
	TargetObjectType  string      `json:"targetObjectType"`
	RequestSpec       interface{} `json:"spec"`
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type RequestPatch PatchBase

type ResponsePatch struct {
	EndpointPath    string           `json:"endpoint_path"`
	PatchOperations []PatchOperation `json:"patch_operations"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

type ResponsePatchList struct {
	PatchList []ResponsePatch `json:"patch_list"`
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
