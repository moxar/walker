package walker

import "fmt"

// PathError is returned when no path can be found between Src and Dst.
type PathError struct {
	Src string
	Dst string
}

// Error implements the error interface.
func (e *PathError) Error() string {
	return fmt.Sprintf("no path between [%s] and [%s]", e.Src, e.Dst)
}
