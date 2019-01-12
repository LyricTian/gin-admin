package context

import "testing"

func TestRouterTitle(t *testing.T) {
	SetRouterTitle("GET", "/api/v1/test/:id/status", "test id")

	title := GetRouterTitle("GET", "/api/v1/test/123/status")
	if title != "test id" {
		t.Error("not the expected value:", title)
	}
}
