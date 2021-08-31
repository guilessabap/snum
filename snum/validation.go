package snum

import "fmt"

type iValidatable interface {
	Validate() error
}

func checkName(name string) error {
	if len(name) < 5 || len(name) > 40 {
		return fmt.Errorf("Length of the name must be between 5 and 40 characters")
	}
	return nil
}
