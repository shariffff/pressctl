package prompt

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
)

// normalizeErr converts huh.ErrUserAborted into a plain "cancelled" error
// so callers don't need to import huh just to check for cancellation.
func normalizeErr(err error) error {
	if errors.Is(err, huh.ErrUserAborted) {
		return fmt.Errorf("cancelled")
	}
	return err
}
