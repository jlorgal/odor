package odor

// Pipeline of filters.
type Pipeline struct {
	filters []Filter
}

// NewPipeline creates a new Pipeline object.
func NewPipeline() *Pipeline {
	return &Pipeline{
		filters: []Filter{},
	}
}

// AddFilters to add a list of filters to the pipeline.
func (p *Pipeline) AddFilters(filters ...Filter) {
	p.filters = append(p.filters, filters...)
}

// HandlePacket handles a packet in the pipeline iterating every filter in sequence.
func (p *Pipeline) HandlePacket(context *Context) FilterAction {
	for _, filter := range p.filters {
		action := filter.Request(context)
		if action == Drop {
			return Drop
		}
	}
	return Accept
}
