package options

// Option is an interface for a single, non-keyed Option, usually as a value for
// an array or embedded map.
type Option struct {
	raw interface{}
}

// NewOption creates an Option wrapper for an interface
func NewOption(raw interface{}) Option {
	return Option{raw}
}

// GetString returns the string value of the option.
// It returns an empty string if the option is not a string.
func (o Option) GetString() string {
	s, ok := o.String()
	if !ok {
		return ""
	}

	return s
}

// String returns the string value of the option and a bool indicating if the option
// is a valid string.
func (o Option) String() (string, bool) {
	s, ok := o.raw.(string)
	if !ok {
		return "", false
	}

	return s, true
}

// GetInt returns the int value of the option.
// It returns -1 if the option is not an int.
func (o Option) GetInt() int {
	i, ok := o.Int()
	if !ok {
		return -1
	}

	return i
}

// Int returns the int value of the option and a bool indicating if the option
// is a valid int.
func (o Option) Int() (int, bool) {
	i, ok := o.raw.(float64)
	if !ok {
		return -1, false
	}

	return int(i), true
}

// GetBool returns the boolean value of the option.
// It returns false if the option is not a boolean.
func (o Option) GetBool() bool {
	b, ok := o.Bool()
	if !ok {
		return false
	}

	return b
}

// Bool returns the boolean value of the option and a bool indicating if the option
// is a valid boolean.
func (o Option) Bool() (bool, bool) {
	b, ok := o.raw.(bool)
	if !ok {
		return false, false
	}

	return b, true
}

// GetArray returns the array value (of options) of the option.
// It returns an empty slice if the option is not an array.
func (o Option) GetArray() []Option {
	a, ok := o.Array()
	if !ok {
		return []Option{}
	}

	return a
}

// Array returns the array value (of options) of the option and a bool indicating
// if the option is a valid array.
func (o Option) Array() ([]Option, bool) {
	var a []Option
	ra, ok := o.raw.([]interface{})
	if !ok {
		return a, false
	}

	for _, op := range ra {
		a = append(a, NewOption(op))
	}

	return a, true
}

// GetMap returns the map value (of strings to options) of the option.
// It returns an empty map if the option is not a map.
func (o Option) GetMap() map[string]Option {
	m, ok := o.Map()
	if !ok {
		return map[string]Option{}
	}

	return m
}

// Map returns the map value (of strings to options) of the option and a bool
// indicating if the option is a valid map.
func (o Option) Map() (map[string]Option, bool) {
	m := make(map[string]Option)
	rm, ok := o.raw.(map[string]interface{})
	if !ok {
		return m, false
	}

	for k, v := range rm {
		m[k] = NewOption(v)
	}

	return m, true
}
