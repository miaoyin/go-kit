package fnutil

// RetryFunc
// @param retryCount
// @param fn(count int) (ok bool, err error)
func RetryFunc(retryCount int, fn func(int) (bool, error)) error {
	var (
		ok  bool
		err error
	)
	for i := 0; i < retryCount; i++ {
		ok, err = fn(i)
		if ok {
			return err
		} else {
			continue
		}
	}
	return err
}
