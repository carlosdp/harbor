package options

// Options is an interface for options passed into chain link configurations.
type Options interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetArray(key string) []Option
	GetMap(key string) map[string]Option
	String(key string) (string, bool)
	Bool(key string) (bool, bool)
	Int(key string) (int, bool)
	Array(key string) ([]Option, bool)
	Map(key string) (map[string]Option, bool)
}

type options struct {
	RawOptions map[string]interface{}
}

// NewOptions takes a raw map from string to interface (from JSON decode) and
// returns a new Options wrapper.
func NewOptions(raw map[string]interface{}) Options {
	return options{raw}
}

func (o options) GetString(key string) string {
	s, ok := o.String(key)
	if !ok {
		return ""
	}

	return s
}

func (o options) String(key string) (string, bool) {
	obj, ok := o.RawOptions[key]
	if !ok {
		return "", false
	}

	s, ok := obj.(string)
	if !ok {
		return "", false
	}

	return s, true
}

func (o options) GetInt(key string) int {
	i, ok := o.Int(key)
	if !ok {
		return -1
	}

	return i
}

func (o options) Int(key string) (int, bool) {
	obj, ok := o.RawOptions[key]
	if !ok {
		return -1, false
	}

	i, ok := obj.(float64)
	if !ok {
		return -1, false
	}

	return int(i), true
}

func (o options) GetBool(key string) bool {
	b, ok := o.Bool(key)
	if !ok {
		return false
	}

	return b
}

func (o options) Bool(key string) (bool, bool) {
	obj, ok := o.RawOptions[key]
	if !ok {
		return false, false
	}

	b, ok := obj.(bool)
	if !ok {
		return false, false
	}

	return b, true
}

func (o options) GetArray(key string) []Option {
	a, ok := o.Array(key)
	if !ok {
		return []Option{}
	}

	return a
}

func (o options) Array(key string) ([]Option, bool) {
	var a []Option
	obj, ok := o.RawOptions[key]
	if !ok {
		return a, false
	}

	ra, ok := obj.([]interface{})
	if !ok {
		return a, false
	}

	for _, op := range ra {
		a = append(a, NewOption(op))
	}

	return a, true
}

func (o options) GetMap(key string) map[string]Option {
	m, ok := o.Map(key)
	if !ok {
		return map[string]Option{}
	}

	return m
}

func (o options) Map(key string) (map[string]Option, bool) {
	m := make(map[string]Option)
	obj, ok := o.RawOptions[key]
	if !ok {
		return m, false
	}

	rm, ok := obj.(map[string]interface{})
	if !ok {
		return m, false
	}

	for k, v := range rm {
		m[k] = NewOption(v)
	}

	return m, true
}
