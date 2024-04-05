package filter

func get_first_el_if_exists(m map[string][]string, key string) (string, bool) {
	var els []string
	var ok bool
	if els, ok = m[key]; !ok {
		return "", false
	}

	if len(els) == 0 {
		return "", false
	}

	return els[0], true
}

func get_val_if_exists(m map[string][]string, key string) ([]string, bool) {
	var els []string
	var ok bool

	if els, ok = m[key]; !ok {
		return nil, false
	}

	if len(els) == 0 {
		return nil, false
	}

	return els, true

}
