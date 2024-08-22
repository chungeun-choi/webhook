package mutating

import (
	"fmt"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigBuilder struct {
	admissionregistrationv1.MutatingWebhookConfiguration
}

func NewMutatingConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (b *ConfigBuilder) WithMetaInfo(name string) *ConfigBuilder {
	b.ObjectMeta.Name = name
	return b
}

func (b *ConfigBuilder) WithWebhook(builder *WebhookConfigBuilder) *ConfigBuilder {
	b.Webhooks = append(b.Webhooks, builder.MutatingWebhook)
	return b
}

type WebhookConfigBuilder struct {
	admissionregistrationv1.MutatingWebhook
}

func NewWebhookConfigBuilder() *WebhookConfigBuilder {
	return &WebhookConfigBuilder{}
}

// WithName sets the name of the webhook
func (b *WebhookConfigBuilder) WithName(name string) *WebhookConfigBuilder {
	b.Name = name
	return b
}

// WithAdmissionReviewVersions sets the admission review versions for the webhook - required
func (b *WebhookConfigBuilder) WithAdmissionReviewVersions(versions ...string) *WebhookConfigBuilder {
	b.AdmissionReviewVersions = versions
	return b
}

// WithSideEffect sets the side effect for the webhook - required
func (b *WebhookConfigBuilder) WithSideEffect(sideEffect admissionregistrationv1.SideEffectClass) *WebhookConfigBuilder {
	b.SideEffects = &sideEffect
	return b
}

func (b *WebhookConfigBuilder) WithFailurePolicy(policy admissionregistrationv1.FailurePolicyType) *WebhookConfigBuilder {
	b.FailurePolicy = &policy
	return b
}

// WithClientConfig sets the client configuration for the webhook - required
func (b *WebhookConfigBuilder) WithClientConfig(url, endpoint string, caByte []byte) *WebhookConfigBuilder {
	//b.ClientConfig.CABundle = caPem

	//Use the direct URL
	url = fmt.Sprintf("%v/%s/%s/%s", url, "patch", endpoint, "trigger")
	b.ClientConfig.URL = &url
	b.ClientConfig.CABundle = caByte

	return b
}

// WithRoles sets the rules for the webhook
func (b *WebhookConfigBuilder) WithRoles(rules ...Rule) *WebhookConfigBuilder {
	for _, rule := range rules {
		var (
			operations []admissionregistrationv1.OperationType
			ruleObj    admissionregistrationv1.Rule
		)

		// Set the rule object
		ruleObj.APIGroups = rule.APIGroups
		ruleObj.APIVersions = rule.APIVersions
		ruleObj.Resources = rule.Resources

		// Convert the string operations to the OperationType
		for _, op := range rule.Operations {
			operations = append(operations, admissionregistrationv1.OperationType(op))
		}

		// Append the rule with operations
		b.Rules = append(b.Rules, admissionregistrationv1.RuleWithOperations{
			Operations: operations,
			Rule:       ruleObj,
		})
	}

	return b
}

func (b *WebhookConfigBuilder) WitNameSpaceSelector(key, value string) *WebhookConfigBuilder {
	// TODO: Implement MatchExpressions
	b.NamespaceSelector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			key: value,
		},
	}
	return b
}

func (b *WebhookConfigBuilder) WithObjectSelector() *WebhookConfigBuilder {
	//TODO: Implement the object selector
	return b
}
