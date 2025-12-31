package googleapi

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestErrors_IsHelpers(t *testing.T) {
	base := errors.New("base")
	if !IsAuthRequiredError(&AuthRequiredError{Service: "gmail", Email: "a@b.com", Cause: base}) {
		t.Fatalf("expected IsAuthRequiredError")
	}
	if !IsRateLimitError(&RateLimitError{RetryAfter: time.Second, Retries: 2}) {
		t.Fatalf("expected IsRateLimitError")
	}
	if !IsCircuitBreakerError(&CircuitBreakerError{}) {
		t.Fatalf("expected IsCircuitBreakerError")
	}
	if !IsQuotaExceededError(&QuotaExceededError{Resource: "gmail"}) {
		t.Fatalf("expected IsQuotaExceededError")
	}
	if !IsNotFoundError(&NotFoundError{Resource: "msg", ID: "id"}) {
		t.Fatalf("expected IsNotFoundError")
	}
	if !IsPermissionDeniedError(&PermissionDeniedError{Resource: "file", Action: "read"}) {
		t.Fatalf("expected IsPermissionDeniedError")
	}
}

func TestErrors_Messages(t *testing.T) {
	base := errors.New("base")
	authErr := &AuthRequiredError{Service: "gmail", Email: "a@b.com", Cause: base}
	if got := authErr.Error(); got != "auth required for gmail a@b.com" {
		t.Fatalf("unexpected: %q", got)
	}
	if !errors.Is(authErr, base) {
		t.Fatalf("expected unwrap to match base")
	}
	if got := (&RateLimitError{RetryAfter: 2 * time.Second, Retries: 3}).Error(); !strings.Contains(got, "retry after 2s") {
		t.Fatalf("unexpected: %q", got)
	}
	if got := (&RateLimitError{Retries: 1}).Error(); !strings.Contains(got, "after 1 retries") {
		t.Fatalf("unexpected: %q", got)
	}
	if got := (&NotFoundError{Resource: "file"}).Error(); got != "file not found" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := (&NotFoundError{Resource: "file", ID: "id1"}).Error(); got != "file not found: id1" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := (&PermissionDeniedError{Resource: "file"}).Error(); got != "permission denied for file" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := (&PermissionDeniedError{Resource: "file", Action: "delete"}).Error(); got != "permission denied: cannot delete file" {
		t.Fatalf("unexpected: %q", got)
	}
}
