package filter

type FilterErrs []error

func (e *FilterErrs) Error() string {
	var errs = ""

	for _, err := range *e {
		errs += ",\n" + err.Error()
	}

	if len(errs) > 0 {
		errs = errs[1:]
	}

	return "[" + errs + "\n]"
}
