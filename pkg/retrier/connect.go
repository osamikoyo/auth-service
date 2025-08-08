package retrier

import "time"

func Connect[T any](retriers uint, connect func() (T, error)) (T, error) {
	timeout := 50 * time.Millisecond

	var (
		value T
		err   error
	)

	for range retriers {
		value, err = connect()
		if err == nil {
			return value, nil
		}

		time.Sleep(timeout)

		timeout *= 2
	}

	return value, err
}
