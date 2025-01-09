package tavily

// OptionManager
type OptionManager struct {
	options map[string]any
}

// NewOptionManager
func NewOptionManager() *OptionManager {
	return &OptionManager{
		options: make(map[string]any),
	}
}

// GetOptionWithDefault
func (o *OptionManager) GetOptionWithDefault(key string, defaultValue any) any {
	if val, ok := o.options[key]; ok {
		return val
	}
	return defaultValue
}

// GetOption
func (o *OptionManager) GetOption(key string) (any, bool) {
	v, ok := o.options[key]
	return v, ok
}

// SetOption
func (o *OptionManager) SetOption(key string, value any) {
	o.options[key] = value
}

// WithOption
type WithOptionHelper func(*OptionManager)

// WithOption
func WithOption(key string, value any) func(*OptionManager) {
	return func(o *OptionManager) {
		o.SetOption(key, value)
	}
}
