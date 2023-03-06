package schema

import "testing"

func TestResourceCreateValidation(t *testing.T) {
	var err error
	var r ResourceCreate
	if err = r.Validate(); err == nil {
		t.Error("should have error")
	}
	r.Code = "test"
	if err = r.Validate(); err == nil {
		t.Error("should have error")
	}
	r.Object = "test"
	if err = r.Validate(); err == nil {
		t.Error("should have error")
	}
	r.Action = "test"
	if err = r.Validate(); err == nil {
		t.Error("should have error")
	}
	r.Status = "disabled"
	if err = r.Validate(); err != nil {
		t.Error("should not have error", err.Error())
	}
}
