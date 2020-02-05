package flags

// Slice is a slice flag type
type Slice []string

func (i *Slice) String() string {
	return ""
}

// Set appends the value into the flag var
func (i *Slice) Set(value string) error {
	*i = append(*i, value)
	return nil
}
