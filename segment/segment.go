package segment

// Segment represents a displayable statusline segment
type Segment interface {
	ID() string
	Render() (string, error)
	Enabled() bool
}
