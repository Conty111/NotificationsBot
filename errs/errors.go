package errs

func CheckError(err error) error {
	if err != nil {
		return err
	}
	return nil
}
