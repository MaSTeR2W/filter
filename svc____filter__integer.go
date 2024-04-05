package filter

import (
	"net/url"
	"strconv"
)

// eq, null, gt, gte, lt, lte

type int_filter struct {
	key       string
	col_alias string
	max       int
	min       int
	check_max bool
	check_min bool
	s_max     string
	s_min     string
}

type IntFilterOpts struct {
	Key       string
	ColAlias  string
	EnableMax bool
	Max       int
	EnableMin bool
	Min       int
}

func NewIntFilter(
	opts IntFilterOpts,
) *int_filter {
	if opts.ColAlias == "" {
		opts.ColAlias = opts.Key
	}

	var f = int_filter{
		key:       opts.Key,
		col_alias: opts.ColAlias,
	}

	if opts.EnableMax {
		f.check_max = true
		f.max = opts.Max
		f.s_max = strconv.Itoa(opts.Max)
	}

	if opts.EnableMin {
		f.check_min = true
		f.min = opts.Min
		f.s_min = strconv.Itoa(opts.Min)
	}

	return &f
}

func (i *int_filter) validate_and_construct(
	v url.Values,
	lang string,
) (string, error) {

	var val string
	var ok bool
	var err error

	if val, ok = get_first_el_if_exists(v, i.key+"[eq]"); ok {
		if val, err = i.validate_val(val, lang); err != nil {
			return "", err
		}

		return i.col_alias + "=" + val, nil
	}

	if _, ok = v[i.key+"[null]"]; ok {
		return i.col_alias + "=NULL", nil
	}

	var conds = []string{}

	if val, ok = get_first_el_if_exists(v, i.key+"[gt]"); ok {
		if val, err = i.validate_val(val, lang); err != nil {
			return "", err
		}

		conds = append(conds, i.col_alias+">"+val)
	} else if val, ok = get_first_el_if_exists(v, i.key+"[gte]"); ok {
		if val, err = i.validate_val(val, lang); err != nil {
			return "", err
		}
		conds = append(conds, i.col_alias+">="+val)
	}

	if val, ok = get_first_el_if_exists(v, i.key+"[lt]"); ok {
		if val, err = i.validate_val(val, lang); err != nil {
			return "", err
		}
		conds = append(conds, i.col_alias+"<"+val)

	} else if val, ok = get_first_el_if_exists(v, i.key+"[lte]"); ok {
		if val, err = i.validate_val(val, lang); err != nil {
			return "", err
		}
		conds = append(conds, i.col_alias+"<="+val)

	}

	var l = len(conds)

	if l == 0 {
		return "", nil
	}

	if l == 1 {
		return conds[0], nil
	}

	return "(" + conds[0] + " AND " + conds[1] + ")", nil
}

func (i *int_filter) validate_val(v string, lang string) (string, error) {
	var num, err = strconv.Atoi(v)

	if err != nil {
		return "", &FilterErr{
			Key:     i.key,
			Value:   v,
			Message: invalid_num_err(lang),
		}
	}

	if i.check_min && num < i.min {
		return "", &FilterErr{
			Key:     i.key,
			Value:   v,
			Message: small_num_err(i.s_min, lang),
		}
	}

	if i.check_max && num > i.max {
		return "", &FilterErr{
			Key:     i.key,
			Value:   v,
			Message: large_num_err(i.s_max, lang),
		}
	}

	return strconv.Itoa(num), nil
}

func invalid_num_err(lang string) string {
	if lang == "ar" {
		return "عدد غير صالح"
	}
	return "invalid number"
}

func small_num_err(s_exp string, lang string) string {
	if lang == "ar" {
		return "يجب أن يكون العدد أكبر من أو يساوي " + s_exp
	}
	return "The number should be greater than or equal to " + s_exp
}

func large_num_err(s_exp string, lang string) string {
	if lang == "ar" {
		return "يجب أن يكون العدد أصغر من أو يساوي " + s_exp
	}

	return "The number should be less than or equal to " + s_exp
}
