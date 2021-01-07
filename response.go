package config

// Response is an example response for docs/tests
type Response struct {
	// ID
	Name string `toml:"name"`

	// Links
	Parent        string `toml:"parent"`
	parentExample *Response

	// Inheritable vals
	Code int               `toml:"code"`
	Data map[string]string `toml:"data"`
}

// InheritFrom ensures family tree is populated and overridden
func (r *Response) InheritFrom(parents map[string]*Response) {
	if parent, ok := parents[r.Parent]; ok {
		if parent.Parent != "" && parent.parentExample == nil {
			// Recursive call to populate family tree from the bottom up
			parent.InheritFrom(parents)
		}

		if r.Code == 0 {
			r.Code = parent.Code
		}

		if r.Data == nil {
			r.Data = parent.Data
		}

		r.parentExample = parent
	}
}
