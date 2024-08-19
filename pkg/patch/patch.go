package patch

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func Patch(name string) *admissionv1.AdmissionResponse {
	var result *admissionv1.AdmissionResponse = new(admissionv1.AdmissionResponse)

	if patchManager, ok := ManagerMap[name]; !ok {
		result.Result = &metav1.Status{
			Message: "No patch manager found",
		}
	} else {
		// Add an annotation to the object
		patchManager.PatchOperations = append(patchManager.PatchOperations, createAnnotationPatch())

		// Marshal the patch operations
		data, err := json.Marshal(patchManager.PatchOperations)
		if err != nil {
			log.Printf("Error: %v", err)

			return result
		}

		// Create the admission response
		result = createAdmissionResponseWithPatch(data)
		if result.Patch == nil {
			result.Result = &metav1.Status{
				Message: "No patch data",
			}
		}

	}

	return result
}

func createAnnotationPatch() PatchOperation {
	return PatchOperation{
		Op:    "add",
		Path:  "metadata/annotations",
		Value: map[string]string{admissionWebhookAnnotationStatusKey: "injected"},
	}
}

func createAdmissionResponseWithPatch(patchData []byte) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchData,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
