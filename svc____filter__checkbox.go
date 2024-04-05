package filter

import (
	"net/url"
	"slices"
	"strconv"
	"strings"
)

func exceed_num_of_available_opts_err(opts_num string, lang string) string {
	if lang == "ar" {
		return "لا يمكن تجاوز عدد الخيارات المتاحة (" + opts_num + ")"
	}
	return "The number of options available (" + opts_num + ") cannot be exceeded"
}

type checkbox_int_filter struct {
	key        string
	col_alias  string
	null_opt   bool
	opts       []int
	s_opts     string
	opts_num   int
	s_opts_num string
}

type CheckboxIntFilterOpts struct {
	Key      string
	ColAlias string
	NullOpt  bool
	Opts     []int
}

func NewCheckboxIntFilter(
	opts CheckboxIntFilterOpts,

) *checkbox_int_filter {
	if opts.ColAlias == "" {
		opts.ColAlias = opts.Key
	}

	var opts_num = len(opts.Opts)

	return &checkbox_int_filter{
		key:       opts.Key,
		col_alias: opts.ColAlias,
		opts:      opts.Opts,
		s_opts: strings.Join(
			sl_of_int_to_sl_of_str[int](opts.Opts),
			", ",
		),
		null_opt:   opts.NullOpt,
		opts_num:   opts_num,
		s_opts_num: strconv.Itoa(opts_num),
	}
}

func (c *checkbox_int_filter) validate_and_construct(
	v url.Values,
	lang string,
) (string, error) {
	var cond string
	var err error

	if cond, err = c.validate_and_construct_vals(v, lang); err != nil {
		return "", err
	}

	if !c.null_opt {
		return cond, nil
	}

	var null string
	var ok bool
	if null, ok = get_first_el_if_exists(v, c.key+"[null]"); ok {

		if null == "0" {
			null = c.col_alias + " IS NOT NULL"
		} else {
			null = c.col_alias + " IS NULL"
		}

		if len(cond) > 0 {
			cond = "(" + cond + " AND " + null + ")"
		} else {
			cond = null
		}
	}

	return cond, nil
}

func (c *checkbox_int_filter) validate_and_construct_vals(
	v url.Values,
	lang string,
) (string, error) {
	var vals []string
	var ok bool

	var op string

	if vals, ok = get_val_if_exists(v, c.key+"[in]"); ok {
		op = " IN "

	} else if vals, ok = get_val_if_exists(v, c.key+"[nin]"); ok {
		op = " NOT IN "
	} else {
		return "", nil
	}

	var err error
	if vals, err = c.validate_arr_of_int(vals, lang); err != nil {
		return "", err
	}

	return c.col_alias + op + "(" + strings.Join(vals, ",") + ")", nil
}

func (c *checkbox_int_filter) validate_arr_of_int(
	v []string,
	lang string,
) ([]string, error) {

	var err error

	if len(v) > c.opts_num {
		return nil, &FilterErr{
			Key:     c.key,
			Value:   v,
			Message: exceed_num_of_available_opts_err(c.s_opts_num, lang),
		}
	}

	var safe_int_str_arr = make([]string, 0, len(v))

	var num int
	for idx, el := range v {

		if num, err = strconv.Atoi(el); err != nil {
			return nil, &FilterErr{
				Key:     c.key,
				Value:   el,
				Path:    []any{idx},
				Message: entry_is_not_num_err(lang),
			}
		}

		if !slices.Contains(c.opts, num) {
			return nil, &FilterErr{
				Key:     c.key,
				Value:   el,
				Path:    []any{idx},
				Message: num_is_not_one_of(c.s_opts, lang),
			}
		}

		// may be strconv.Atoi has unexpected behaviour
		// so cannot trust the original string value
		safe_int_str_arr = append(safe_int_str_arr, strconv.Itoa(num))
	}

	return safe_int_str_arr, nil
}

func entry_is_not_num_err(lang string) string {
	if lang == "ar" {
		return "هذا ليس عددا"
	}
	return "This is not a number"
}

func num_is_not_one_of(one_of string, lang string) string {
	if lang == "ar" {
		return "يجب أن يكون العدد واحد من: (" + one_of + ")"
	}
	return "The number should be one of: (" + one_of + ")"
}

type checkbox_str_filter struct {
	key        string
	col_alias  string
	null_opt   bool
	opts       []string
	s_opts     string
	opts_num   int
	s_opts_num string
}

type CheckboxStrFilterOpts struct {
	Key      string
	ColAlias string
	NullOpt  bool
	Opts     []string
}

func NewCheckboxStrFilter(
	opts CheckboxStrFilterOpts,
) *checkbox_str_filter {
	if opts.ColAlias == "" {
		opts.ColAlias = opts.Key
	}

	var opts_num = len(opts.Opts)

	var s_opts = ""
	for _, opt := range opts.Opts {
		s_opts += ", " + opt
	}

	if len(s_opts) > 0 {
		s_opts = s_opts[2:]
	}

	return &checkbox_str_filter{
		key:        opts.Key,
		col_alias:  opts.ColAlias,
		null_opt:   opts.NullOpt,
		opts:       opts.Opts,
		s_opts:     s_opts,
		opts_num:   opts_num,
		s_opts_num: strconv.Itoa(opts_num),
	}
}

func (c *checkbox_str_filter) validate_and_construct(
	v url.Values,
	lang string,
) (string, error) {
	var cond string
	var err error

	if cond, err = c.validate_and_construct_vals(v, lang); err != nil {
		return "", err
	}

	if !c.null_opt {
		return cond, nil
	}

	var null string
	var ok bool
	if null, ok = get_first_el_if_exists(v, c.key+"[null]"); ok {

		if null == "0" {
			null = c.col_alias + " IS NOT NULL"
		} else {
			null = c.col_alias + " IS NULL"
		}

		if len(cond) > 0 {
			cond = "(" + cond + " AND " + null + ")"
		} else {
			cond = null
		}
	}

	return cond, nil
}

func (c *checkbox_str_filter) validate_and_construct_vals(
	v url.Values,
	lang string,
) (string, error) {
	var vals []string
	var ok bool

	var op string

	if vals, ok = get_val_if_exists(v, c.key+"[in]"); ok {
		op = " IN "

	} else if vals, ok = get_val_if_exists(v, c.key+"[nin]"); ok {
		op = " NOT IN "
	} else {
		return "", nil
	}

	var err error
	if vals, err = c.validate_arr_of_str(vals, lang); err != nil {
		return "", err
	}

	return c.col_alias + op + "(" + strings.Join(vals, ",") + ")", nil
}

func (c *checkbox_str_filter) validate_arr_of_str(
	v []string,
	lang string,
) ([]string, error) {

	var escaped_strs = make([]string, 0, len(v))

	for idx, el := range v {

		if !slices.Contains(c.opts, el) {
			return nil, &FilterErr{
				Key:     c.key,
				Value:   el,
				Path:    []any{idx},
				Message: str_is_not_one_of(c.s_opts, lang),
			}
		}

		escaped_strs = append(escaped_strs, to_escaped_string(el))

	}

	return escaped_strs, nil
}

func str_is_not_one_of(one_of string, lang string) string {
	if lang == "ar" {
		return "يجب أن يكون الخيار واحد من: (" + one_of + ")"
	}
	return "The option should be one of: (" + one_of + ")"
}
