package filter

import (
	"net/url"
	"strconv"
)

type str_filter struct {
	key           string
	col_alias     string
	check_max_len bool
	max_len       int
	s_max_len     string
}

type StrFilterOpts struct {
	Key          string
	ColAlias     string
	EnableMaxLen bool
	MaxLen       int
}

func NewStrFilter(opts StrFilterOpts) *str_filter {
	if opts.ColAlias == "" {
		opts.ColAlias = opts.Key
	}

	var f = str_filter{
		key:       opts.Key,
		col_alias: opts.ColAlias,
	}

	if opts.EnableMaxLen {
		f.check_max_len = true
		f.max_len = opts.MaxLen
		f.s_max_len = strconv.Itoa(opts.MaxLen)
	}

	return &f
}

func (s *str_filter) validate_and_construct(
	v url.Values,
	lang string,
) (string, error) {

	var val string
	var ok bool
	var err error

	if val, ok = get_first_el_if_exists(v, s.key+"[eq]"); ok {
		if err = s.does_val_exceed_max_len(val, lang); err != nil {
			return "", err
		}

		return s.col_alias + "=" + to_escaped_string(val), nil
	}

	if val, ok = get_first_el_if_exists(v, s.key+"[null]"); ok {
		if err = s.does_val_exceed_max_len(val, lang); err != nil {
			return "", err
		}
		return s.col_alias + "=NULL", nil
	}

	if val, ok = get_first_el_if_exists(v, s.key+"[sw]"); ok {
		if err = s.does_val_exceed_max_len(val, lang); err != nil {
			return "", err
		}
		return s.col_alias + " LIKE '" + escape_single_quote(val) + "%'", nil
	}

	if val, ok = get_first_el_if_exists(v, s.key+"[ew]"); ok {
		if err = s.does_val_exceed_max_len(val, lang); err != nil {
			return "", err
		}
		return s.col_alias + " LIKE '%" + escape_single_quote(val) + "'", nil
	}

	if val, ok = get_first_el_if_exists(v, s.key+"[ct]"); ok {
		if err = s.does_val_exceed_max_len(val, lang); err != nil {
			return "", err
		}
		return s.col_alias + " LIKE '%" + escape_single_quote(val) + "%'", nil
	}

	return "", nil
}

func (s *str_filter) does_val_exceed_max_len(v string, lang string) error {
	if !s.check_max_len {
		return nil
	}
	var l = len(v)
	if l > s.max_len {
		return &FilterErr{
			Key:     s.key,
			Value:   v,
			Message: long_str_err(s.s_max_len, l, lang),
		}
	}
	return nil
}

func long_str_err(s_exp string, got int, lang string) string {
	var s_got = strconv.Itoa(got)
	if lang == "ar" {
		return "يجب تقصير هذا النص إلى " + s_exp + " من الحروف أو أقل (أنت حاليا تستخدم " + s_got + " من الحروف)"
	}

	return "Should shorten this text to " + s_exp + " characters (you are currently using " + s_got + " characters)"
}
