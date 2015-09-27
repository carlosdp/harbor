package options_test

import (
	"testing"

	"github.com/carlosdp/harbor/options"
)

func TestWrapsString(t *testing.T) {
	t.Parallel()
	mi := map[string]interface{}{
		"test": "success",
	}

	ops := options.NewOptions(mi)
	if ops.GetString("test") != "success" {
		t.Fatal("should have wrapped string correctly")
	}
}

func TestWrapsInt(t *testing.T) {
	t.Parallel()
	mi := map[string]interface{}{
		"test": 5.0,
	}

	ops := options.NewOptions(mi)
	if ops.GetInt("test") != 5 {
		t.Fatal("should have wrapped int correctly")
	}
}

func TestWrapsBool(t *testing.T) {
	t.Parallel()
	mi := map[string]interface{}{
		"test": true,
	}

	ops := options.NewOptions(mi)
	if !ops.GetBool("test") {
		t.Fatal("should have wrapped bool correctly")
	}
}

func TestWrapsArray(t *testing.T) {
	t.Parallel()
	mi := map[string]interface{}{
		"test": []interface{}{"success"},
	}

	ops := options.NewOptions(mi)
	a := ops.GetArray("test")
	if len(a) != 1 || a[0].GetString() != "success" {
		t.Fatal("should have wrapped array correctly")
	}
}

func TestWrapsMap(t *testing.T) {
	t.Parallel()
	mi := map[string]interface{}{
		"test": map[string]interface{}{"test": "success"},
	}

	ops := options.NewOptions(mi)
	m := ops.GetMap("test")
	success, ok := m["test"]
	if !ok || success.GetString() != "success" {
		t.Fatal("should have wrapped map correctly")
	}
}
