package mutating

import admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

type Rule struct {
	APIGroups   []string `json:"apiGroup"`
	APIVersions []string `json:"apiVersion"`
	Resources   []string `json:"resource"`
	Operations  []string `json:"operations"`
}

type RequestAddRuleBody struct {
	Rule Rule `json:"rule"`
}

type RequestAddRulesBody struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

type ResponseBody struct {
	Message string `json:"message"`
}

type ResponseAddRulesBody struct {
	ResponseBody
	ConfigBuilder
}

type ResponseGetRulesBody struct {
	Message              string                                               `json:"message"`
	WebhookConfiguration admissionregistrationv1.MutatingWebhookConfiguration `json:"webhookConfiguration"`
}
