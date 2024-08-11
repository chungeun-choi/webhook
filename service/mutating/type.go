package mutating

type Rule struct {
	APIGroups   []string `json:"APIGroup"`
	APIVersions []string `json:"APIVersion"`
	Resources   []string `json:"Resource"`
	Operations  []string `json:"Operations"`
}

type RequestAddRuleBody struct {
	Rule Rule `json:"Rule"`
}

type RequestAddRulesBody struct {
	Rules []Rule `json:"Rules"`
}

type ResponseBody struct {
	Message string `json:"Message"`
}
