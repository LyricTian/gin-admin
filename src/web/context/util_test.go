package context

import "testing"

func TestRouterTitle(t *testing.T) {
	SetRouterTitle("GET", "/api/v1/test/:id/status", "test id")

	title, key := GetRouterTitleAndKey("GET", "/api/v1/test/123/status")
	if title != "test id" {
		t.Error("not the expected value:", title)
	}

	if key != "GET-/api/v1/test/:id/status" {
		t.Error("not the expected value:", key)
	}
}
