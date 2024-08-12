package patch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPatchManager(t *testing.T) {
	endpoint := "/example"
	pm := NewPatchManager(endpoint)

	if pm == nil {
		t.Errorf("Expected PatchManager, got nil")
	}

	if pm.EndpointPath != endpoint {
		t.Errorf("Expected endpoint %s, got %s", endpoint, pm.EndpointPath)
	}

	if ManagerMap[endpoint] != pm {
		t.Errorf("Expected PatchManager in ManagerMap")
	}
}

func TestUpdatePatchOperation(t *testing.T) {
	endpoint := "/example"
	NewPatchManager(endpoint)

	req := RequestPatch{
		Objects: []Object{
			{
				Op:                "add",
				RequestObjectType: "container",
				TargetObjectType:  "deployment",
				RequestSpec:       "example-value",
			},
		},
	}

	pm := NewPatchManager(endpoint)
	_ = pm.AddPatchOperations(req)

	err := pm.updatePatchOperation(req)
	assert.NoError(t, err)

	pm, exists := ManagerMap[endpoint]
	assert.True(t, exists)
	assert.Len(t, pm.GetPatchOperations(), 1)
}

func TestClearPatchOperations(t *testing.T) {
	endpoint := "/example"
	pm := NewPatchManager(endpoint)

	_ = pm.AddPatchOperations(RequestPatch{
		Objects: []Object{
			{
				Op:                "add",
				RequestObjectType: "container",
				TargetObjectType:  "deployment",
				RequestSpec:       "example-value",
			},
		},
	})

	pm.ClearPatchOperations()
	assert.Empty(t, pm.GetPatchOperations())
}

func TestAddPatchOperations(t *testing.T) {
	endpoint := "/example"
	pm := NewPatchManager(endpoint)

	req := RequestPatch{
		Objects: []Object{
			{
				Op:                "add",
				RequestObjectType: "container",
				TargetObjectType:  "deployment",
				RequestSpec:       "example-value",
			},
		},
	}

	err := pm.AddPatchOperations(req)
	assert.NoError(t, err)
	assert.Len(t, pm.GetPatchOperations(), 1)
}

func TestBasePath(t *testing.T) {
	testCases := []struct {
		requestObjectType string
		targetObjectType  string
		expectedPath      string
		expectError       bool
	}{
		{"container", "deployment", "/spec/template/spec/containers", false},
		{"pod", "deployment", "/spec/template/spec", false},
		{"volume", "deployment", "/spec/template/spec/volumes", false},
		{"unsupported", "deployment", "", true},
		{"container", "unsupported", "", true},
	}

	for _, tc := range testCases {
		path, err := basePath(tc.requestObjectType, tc.targetObjectType)
		if tc.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPath, path)
		}
	}
}
