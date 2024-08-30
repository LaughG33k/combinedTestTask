package pkg

import "time"

func Rerty(fn func() error, attempt int, timeSleep time.Duration) (err error) {

	for attempt > 0 {

		if err = fn(); err != nil {

			time.Sleep(timeSleep)
			attempt--

			continue
		}

		return nil

	}

	return err

}
