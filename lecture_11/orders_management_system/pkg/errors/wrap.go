package errors

import "fmt"

// Wrap - is sugar for fmt.Errorf("%s: %w")
func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}
