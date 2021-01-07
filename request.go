package config

// Request is an example requests for docs/tests
type Request struct {
	// ID
	Name  string `toml:"name"`
	Group string `toml:"group"`

	Query map[string]string `toml:"query"`
	Body  map[string]string `toml:"body"`

	Responses        []string    `toml:"responses"`
	ResponseExamples []*Response `toml:"-"`

	// Links
	Parent        string `toml:"parent"`
	parentExample *Request
}

// InheritFrom ensures family tree is populated and overridden
func (r *Request) InheritFrom(parents map[string]*Request) {
	if parent, ok := parents[r.Parent]; ok {
		if parent.Parent != "" && parent.parentExample == nil {
			// Recursive call to populate family tree from the bottom up
			parent.InheritFrom(parents)
		}

		if r.Query == nil {
			r.Query = parent.Query
		}

		if r.Body == nil {
			r.Body = parent.Body
		}

		if len(parent.Responses) > 0 {
			r.Responses = append(r.Responses, parent.Responses...)

			if len(parent.ResponseExamples) == 0 {
				r.ResponseExamples = append(r.ResponseExamples, parent.ResponseExamples...)
			}
		}
	}
}
