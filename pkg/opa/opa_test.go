package opa_test

import (
	"context"
	"testing"

	"github.com/gocastsian/roham/pkg/opa"
	"github.com/stretchr/testify/assert"
)

func TestNewOPAEvaluator(t *testing.T) {
	policy := `package test

default allow = false

allow if{
	input.user == "admin"
}`

	cfg := opa.Config{
		Package: "test",
		Rule:    "allow",
		Policy:  policy,
		IsPath:  false,
	}

	evaluator, err := opa.NewOPAEvaluator(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, evaluator)
}

func TestEvaluate_Success(t *testing.T) {
	policy := `package test

default allow = false

allow if{
	input.user == "admin"
}`

	cfg := opa.Config{
		Package: "test",
		Rule:    "allow",
		Policy:  policy,
		IsPath:  false,
	}

	evaluator, err := opa.NewOPAEvaluator(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, evaluator)

	ctx := context.Background()
	input := map[string]interface{}{"user": "admin"}
	err = evaluator.Evaluate(ctx, input)
	assert.NoError(t, err)
}

func TestEvaluate_Failure(t *testing.T) {
	policy := `package test

default allow = false

allow if{
	input.user == "admin"
}`

	cfg := opa.Config{
		Package: "test",
		Rule:    "allow",
		Policy:  policy,
		IsPath:  false,
	}

	evaluator, err := opa.NewOPAEvaluator(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, evaluator)

	ctx := context.Background()
	input := map[string]interface{}{"user": "guest"}
	err = evaluator.Evaluate(ctx, input)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "bindings results[[{[true] map[x:false]}]] ok[true]")
}
