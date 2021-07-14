// Package testutils contains all mocks, assertions, helpers, … used for testing the libPackage
package testutils

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

var AssertMessageDefault = "%s has wrong %s: target: %q, actual: %q"

type Labels map[string]string

func AssertLabel(t *testing.T, labelsActual Labels, keyTarget, valueTarget, containsMsg string) {
	if assert.Containsf(
		t,
		labelsActual,
		keyTarget,
		containsMsg,
		keyTarget,
	) {
		assert.Equalf(
			t,
			labelsActual[keyTarget],
			valueTarget,
			"value of label %q should be %q, got: %q",
			keyTarget,
			valueTarget,
			labelsActual[keyTarget],
		)
	}
}

func AssertAnnotation(t *testing.T, annotationsActual map[string]string, keyTarget string, valueTarget pulumi.StringInput, containsMsg string) {
	if assert.Containsf(
		t,
		annotationsActual,
		keyTarget,
		containsMsg,
		keyTarget,
	) {
		assert.Equalf(
			t,
			pulumi.String(annotationsActual[keyTarget]),
			valueTarget,
			"value of annotation %q should be %q, got: %q",
			keyTarget,
			valueTarget,
			annotationsActual[keyTarget],
		)
	}
}

func Equalf(t *testing.T, resourceName, subjectName string, actual, target interface{}) {
	assert.Equalf(t, actual, target, AssertMessageDefault, resourceName, subjectName, target, actual)
}

func Containsf(t *testing.T, resourceName, subjectName string, actual, target interface{}) {
	assert.Containsf(t, actual, target, AssertMessageDefault, resourceName, subjectName, target, actual)
}
