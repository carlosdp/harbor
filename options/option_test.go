package options_test

import (
	"testing"

	"github.com/carlosdp/harbor/options"
)

func TestOptionWrapsString(t *testing.T) {
	t.Parallel()
	i := interface{}("success")

	op := options.NewOption(i)
	if op.GetString() != "success" {
		t.Fatal("should have wrapped string correctly")
	}
}

func TestOptionWrapsInt(t *testing.T) {
	t.Parallel()
	i := interface{}(5.0)

	op := options.NewOption(i)
	if op.GetInt() != 5 {
		t.Fatal("should have wrapped int correctly")
	}
}

func TestOptionWrapsBool(t *testing.T) {
	t.Parallel()
	i := interface{}(true)

	op := options.NewOption(i)
	if !op.GetBool() {
		t.Fatal("should have wrapped bool correctly")
	}
}

func TestOptionWrapsArray(t *testing.T) {
	t.Parallel()
	i := []interface{}{"success"}

	op := options.NewOption(i)
	if len(op.GetArray()) != 1 || op.GetArray()[0].GetString() != "success" {
		t.Fatal("should have wrapped array correctly")
	}
}

func TestOptionWrapsMap(t *testing.T) {
	t.Parallel()
	i := map[string]interface{}{
		"test": "success",
	}

	op := options.NewOption(i)
	m := op.GetMap()
	success, ok := m["test"]
	if !ok || success.GetString() != "success" {
		t.Fatal("should have wrapped map correctly")
	}
}
