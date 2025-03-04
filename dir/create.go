package dir

import "os"

func Create(p string) error {
	if err := os.MkdirAll(p, 0755); err != nil {
		return err
	}

	return nil
}
