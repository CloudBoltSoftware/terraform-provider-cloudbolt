package cmp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func TestWithPanicRecovery_RecoversAndAddsDiagnostic(t *testing.T) {
	var diags diag.Diagnostics

	func() {
		defer withPanicRecovery(&diags, "Create")

		panic("boom")
	}()

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}

	d := diags[0]

	if d.Severity != diag.Error {
		t.Errorf("expected Severity=Error, got %v", d.Severity)
	}

	if d.Summary != "CloudBolt provider crashed" {
		t.Errorf("unexpected Summary: %q", d.Summary)
	}

	expectedDetail := "panic during Create: boom"
	if d.Detail != expectedDetail {
		t.Errorf("expected Detail %q, got %q", expectedDetail, d.Detail)
	}
}


func TestWithPanicRecovery_NoPanic_NoDiagnostic(t *testing.T) {
	var diags diag.Diagnostics

	func() {
		defer withPanicRecovery(&diags, "Read")

		// no panic
	}()

	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d", len(diags))
	}
}

