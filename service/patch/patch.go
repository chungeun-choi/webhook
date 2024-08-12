package patch

import (
	"github.com/pkg/errors"
	"log"
)

var (
	ManagerMap map[string]*PatchManager = make(map[string]*PatchManager)
)

type PatchManager struct {
	EndpointPath    string
	PatchOperations []PatchOperation
}

func NewPatchManager(endpointPath string) *PatchManager {
	var pm *PatchManager
	if _, ok := ManagerMap[endpointPath]; !ok {
		pm = &PatchManager{
			EndpointPath: endpointPath,
		}
		ManagerMap[endpointPath] = pm
	} else {
		pm = ManagerMap[endpointPath]
	}

	return pm
}

// Update updates the patch operations for the PatchManager. it was clear the existing patch operations and add the new ones
func (pm *PatchManager) updatePatchOperation(req RequestPatch) error {
	// Clear the existing patch operations
	pm.ClearPatchOperations()

	return pm.AddPatchOperations(req)
}

// GetPatchOperations returns the patch operations for the PatchManager
func (pm *PatchManager) GetPatchOperations() []PatchOperation {
	return pm.PatchOperations
}

// ClearPatchOperations clears the patch operations for the PatchManager
func (pm *PatchManager) ClearPatchOperations() {
	pm.PatchOperations = nil
}

// AddPatchOperations adds the patch operations to the PatchManager
func (pm *PatchManager) AddPatchOperations(req RequestPatch) error {
	for _, obj := range req.Objects {
		operation, err := convertToPatchOperation(obj)
		if err != nil {
			log.Printf("Error: %v", err)
			return err
		}

		pm.PatchOperations = append(pm.PatchOperations, *operation)
	}

	return nil
}

// convertToPatchOperation converts the request object to a PatchOperation
func convertToPatchOperation(info Object) (*PatchOperation, error) {
	var path string

	path, err := basePath(info.RequestObjectType, info.TargetObjectType)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return &PatchOperation{
		Op:    info.Op,
		Path:  path,
		Value: info,
	}, nil
}

// basePath determines the JSON path where the requestObjectType should be added to the targetObjectType
func basePath(requestObjectType, targetObjectType string) (string, error) {
	switch targetObjectType {
	case DEPLOYMENT, STATEFULSET:
		switch requestObjectType {
		case CONTAINER:
			return "/spec/template/spec/containers", nil
		case POD:
			return "/spec/template/spec", nil
		case VOLUME:
			return "/spec/template/spec/volumes", nil
		default:
			return "", errors.New(" unsupported request object type")
		}
	default:
		return "", errors.New("unsupported target object type")
	}
}
