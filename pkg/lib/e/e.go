package e

import "fmt"

func Warp(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
