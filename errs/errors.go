package errs

import "log"

func CheckError(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
