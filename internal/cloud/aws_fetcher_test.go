package cloud_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/driftcheck/internal/cloud"
)

// stubFetcher satisfies a simple interface for unit testing without real AWS.
type stubFetcher struct {
	attrs cloud.ResourceAttributes
	err   error
}

func (s *stubFetcher) FetchResource(_ context.Context, _, _ string) (cloud.ResourceAttributes, error) {
	return s.attrs, s.err
}

func TestFetchResource_UnsupportedType(t *testing.T) {
	ctx := context.Background()
	// We can't call NewFetcher without real AWS creds, so we test the
	// dispatch logic indirectly via a minimal integration shim.
	_ = ctx
	// Validate that an unsupported type surfaces the right error message.
	f := &cloud.Fetcher{}
	_, err := f.FetchResource(ctx, "aws_lambda_function", "my-fn")
	if err == nil {
		t.Fatal("expected error for unsupported resource type, got nil")
	}
	expected := "unsupported resource type: aws_lambda_function"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestStubFetcher_ReturnsAttrs(t *testing.T) {
	stub := &stubFetcher{
		attrs: cloud.ResourceAttributes{"instance_type": "t3.micro"},
	}
	attrs, err := stub.FetchResource(context.Background(), "aws_instance", "i-abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["instance_type"] != "t3.micro" {
		t.Errorf("instance_type = %v, want t3.micro", attrs["instance_type"])
	}
}

func TestStubFetcher_PropagatesError(t *testing.T) {
	expectedErr := errors.New("resource not found")
	stub := &stubFetcher{err: expectedErr}
	_, err := stub.FetchResource(context.Background(), "aws_instance", "i-missing")
	if !errors.Is(err, expectedErr) {
		t.Errorf("err = %v, want %v", err, expectedErr)
	}
}
