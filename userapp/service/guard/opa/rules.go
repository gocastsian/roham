package opa

import (
	_ "embed"
)

// These are the current set of rules we have for auth.
const (
	RuleCheckRequestOnly = "rule_check_request_only"
)

// Package name of our rego code.
const (
	Package string = "roham.rego"
)

// Core OPA policies.
var (
	//go:embed rego/authorization.rego
	RegoAuthorization string
)
