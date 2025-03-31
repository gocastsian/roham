package opa

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/v1/rego"
)

type Config struct {
	Package string `koanf:"package"`
	Rule    string `koanf:"rule"`
	Policy  string `koanf:"policy"`  // either content or file path
	IsPath  bool   `koanf:"is_path"` // true if Policy is a file path
}

type OPAEvaluator struct {
	config Config
	query  rego.PreparedEvalQuery
}

// NewOPAEvaluator initializes OPA with either inline policy or policy from file
func NewOPAEvaluator(cfg Config) (*OPAEvaluator, error) {
	var policyContent string
	var err error

	if cfg.IsPath {
		policyContent, err = readPolicyFile(cfg.Policy)
		if err != nil {
			return nil, fmt.Errorf("failed to read policy file: %w", err)
		}
	} else {
		policyContent = cfg.Policy
	}

	query := fmt.Sprintf("x = data.%s.%s", cfg.Package, cfg.Rule)

	r := rego.New(
		rego.Query(query),
		rego.Module("policy.rego", policyContent),
	)

	q, err := r.PrepareForEval(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to prepare rego query: %w", err)
	}

	return &OPAEvaluator{query: q}, nil
}

func readPolicyFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Evaluate executes the policy with given input data.
func (o *OPAEvaluator) Evaluate(ctx context.Context, input map[string]interface{}) error {
	results, err := o.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !ok || !result {
		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
	}

	return nil
}
