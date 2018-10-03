package scraper

// State interface for all states
type State interface {
	Do() State
}

// NullObject for the state to exit Context Run
var exitState State

// Context to run states
type Context struct {
	current State
}

// NewContext creates new Context
func NewContext(s State) Context {
	ret := Context{s}
	return ret
}

// Run starts the execution of the Context
func (c *Context) Run() {
	for c.current != exitState {
		c.current = c.current.Do()
	}
}
