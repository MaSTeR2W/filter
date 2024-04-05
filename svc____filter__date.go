package filter

import (
	"net/url"
	"strings"
	"time"
)

const DATE_LAYOUT = "2006-01-02"

type date_filter struct {
	key          string
	col_alias    string
	after_now    bool
	check_after  bool
	after_unix   int64
	after_date   string
	before_now   bool
	check_before bool
	before_unix  int64
	before_date  string
	null_opt     bool
}

type DateFilterOpts struct {
	Key       string
	ColAlias  string
	AfterNow  bool
	After     string // should be in (YYYY-MM-DD) format
	BeforeNow bool
	Before    string // should be in (YYYY-MM-DD) format
	NullOpt   bool
}

func MustCreateNewDateFilter(opts DateFilterOpts) *date_filter {
	var ft, err = NewDateFilter(opts)

	if err != nil {
		panic(err)
	}
	return ft
}

func NewDateFilter(opts DateFilterOpts) (*date_filter, error) {

	if opts.ColAlias == "" {
		opts.ColAlias = opts.Key
	}

	var d_filter = date_filter{
		key:        opts.Key,
		col_alias:  opts.ColAlias,
		after_now:  opts.AfterNow,
		before_now: opts.BeforeNow,
		null_opt:   opts.NullOpt,
	}

	var err error

	if !opts.AfterNow && opts.After != "" {
		var t_after time.Time
		if t_after, err = time.Parse(DATE_LAYOUT, opts.After); err != nil {
			return nil, err
		}
		d_filter.after_unix = t_after.Unix()
		d_filter.after_date = t_after.Format(DATE_LAYOUT)
		d_filter.check_after = true
	}

	if !opts.BeforeNow && opts.Before != "" {
		var t_before time.Time

		if t_before, err = time.Parse(DATE_LAYOUT, opts.Before); err != nil {
			return nil, err
		}

		d_filter.before_unix = t_before.Unix()
		d_filter.before_date = t_before.Format(DATE_LAYOUT)
		d_filter.check_before = true
	}

	return &d_filter, nil
}

func (d *date_filter) validate_and_construct(
	v url.Values,
	lang string,
) (string, error) {
	// pr, pre, ps, pse, eq, null

	var input string
	var ok bool
	var err error
	var cond string

	if input, ok = get_first_el_if_exists(v, d.key+"[eq]"); ok {
		if input, err = d.validate_val(input, "eq", lang); err != nil {
			return "", err
		}
		cond = d.col_alias + "=" + wrap_with_single_quote(input)
	} else {
		var conds = []string{}
		if input, ok = get_first_el_if_exists(v, d.key+"[pr]"); ok {

			if input, err = d.validate_val(input, "pr", lang); err != nil {
				return "", err
			}
			conds = append(conds, d.col_alias+"<"+wrap_with_single_quote(input))

		} else if input, ok = get_first_el_if_exists(v, d.key+"[pre]"); ok {

			if input, err = d.validate_val(input, "pre", lang); err != nil {
				return "", err
			}
			conds = append(conds, d.col_alias+"<="+wrap_with_single_quote(input))
		}

		if input, ok = get_first_el_if_exists(v, d.key+"[ps]"); ok {

			if input, err = d.validate_val(input, "ps", lang); err != nil {
				return "", err
			}
			conds = append(conds, d.col_alias+">"+wrap_with_single_quote(input))

		} else if input, ok = get_first_el_if_exists(v, d.key+"[pse]"); ok {

			if input, err = d.validate_val(input, "pse", lang); err != nil {
				return "", err
			}
			conds = append(conds, d.col_alias+">="+wrap_with_single_quote(input))
		}

		switch len(conds) {
		case 1:
			cond = conds[0]
		case 2:
			cond = "(" + strings.Join(conds, " AND ") + ")"
		}
	}

	if d.null_opt {
		if input, ok = get_first_el_if_exists(v, d.key+"[null]"); ok {
			if cond != "" {
				if input != "0" {
					cond = "(" + cond + " OR " + d.col_alias + " IS NULL)"
				}
			} else {
				if input == "0" {
					cond = d.col_alias + " IS NOT NULL"
				} else {
					cond = d.col_alias + " IS NULL"
				}
			}
		}
	}

	return cond, nil
}

func (d *date_filter) validate_val(input string, op string, lang string) (string, error) {

	var t, err = time.Parse(DATE_LAYOUT, input)

	if err != nil {
		return "", &FilterErr{
			Key:     d.key,
			Value:   input,
			Path:    []any{op},
			Message: invalid_date_err(lang),
		}
	}

	input = t.Format(DATE_LAYOUT)

	var input_unix = t.Unix()

	if d.after_now || d.before_now {
		var now_t = time.Now()

		var now_unix = now_t.Unix()

		if d.after_now && input_unix < now_unix {
			return "", &FilterErr{
				Key:     d.key,
				Value:   input,
				Path:    []any{op},
				Message: early_date_err(now_t.Format(DATE_LAYOUT), input, lang),
			}
		}

		if d.before_now && input_unix > now_unix {
			return "", &FilterErr{
				Key:     d.key,
				Value:   input,
				Path:    []any{op},
				Message: late_date_err(now_t.Format(DATE_LAYOUT), input, lang),
			}
		}
	}

	if d.check_after && input_unix < d.after_unix {
		return "", &FilterErr{
			Key:     d.key,
			Value:   input,
			Path:    []any{op},
			Message: early_date_err(d.after_date, input, lang),
		}
	}

	if d.check_before && input_unix > d.before_unix {
		return "", &FilterErr{
			Key:     d.key,
			Value:   input,
			Path:    []any{op},
			Message: late_date_err(d.before_date, input, lang),
		}
	}

	return input, nil

}

func invalid_date_err(lang string) string {
	if lang == "ar" {
		return "التاريخ غير صالح"
	}
	return "The date is invalid"
}

func early_date_err(after string, input string, lang string) string {
	if lang == "ar" {
		return "يجب أن يكون التاريخ بعد (" + after + "), التاريخ الذي أدخلته (" + input + ")"
	}
	return "The date should be after (" + after + "), the date you entered (" + input + ")"
}

func late_date_err(before string, input string, lang string) string {
	if lang == "ar" {
		return "يجب أن يكون التاريخ قبل (" + before + "), التاريخ الذي أدخلته (" + input + ")"
	}
	return "The date should be before (" + before + "), the date you entered (" + input + ")"
}
