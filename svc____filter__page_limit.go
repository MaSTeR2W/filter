package filter

import (
	"net/url"
	"strconv"
)

type paginator struct {
	limit_min        int
	s_limit_min      string
	enable_limit_max bool
	limit_max        int
	s_limit_max      string
}

type paginator_opts struct {
	limit_min        int
	enable_limit_max bool
	limit_max        int
}

func new_pagintor(opts paginator_opts) *paginator {
	var p = paginator{
		limit_min:   opts.limit_min,
		s_limit_min: strconv.Itoa(opts.limit_min),
	}

	if opts.enable_limit_max {
		p.enable_limit_max = true
		p.limit_max = opts.limit_max
		p.s_limit_max = strconv.Itoa(opts.limit_max)
	}
	return &p
}

func (p *paginator) paginate(v url.Values, lang string) (string, error) {
	var page, limit, ok, err = get_limit_page(v, lang)

	if err != nil {
		return "", err
	}

	if !ok {
		return "", nil
	}

	if limit < p.limit_min {
		return "", &FilterErr{
			Key:     "$limit",
			Value:   limit,
			Message: limit_min_err(p.s_limit_min, lang),
		}
	}

	if p.enable_limit_max && limit > p.limit_max {
		return "", &FilterErr{
			Key:     "$limit",
			Value:   limit,
			Message: limit_max_err(p.s_limit_max, lang),
		}
	}

	return "LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(limit*(page-1)), nil
}

func limit_min_err(s_exp string, lang string) string {
	if lang == "ar" {
		return "يجب أن يكون الحد " + s_exp + " على الأقل"
	}
	return "The limit should be at least " + s_exp
}

func limit_max_err(s_exp string, lang string) string {
	if lang == "ar" {
		return "يجب ألا يتجاوز الحد " + s_exp
	}
	return "The limit should not exceed " + s_exp
}

func get_limit_page(v url.Values, lang string) (int, int, bool, error) {
	var (
		page    int
		s_page  string
		limit   int
		ok_page bool
		err     error
	)

	if s_page, ok_page = get_first_el_if_exists(v, "$page"); ok_page {
		if page, err = strconv.Atoi(s_page); err != nil {
			return 0, 0, false, &FilterErr{
				Key:     "$page",
				Value:   s_page,
				Message: invalid_page_num_err(lang),
			}
		}
	}

	var (
		s_limit  string
		ok_limit bool
	)

	if s_limit, ok_limit = get_first_el_if_exists(v, "$limit"); ok_limit {
		if limit, err = strconv.Atoi(s_limit); err != nil {
			return 0, 0, false, &FilterErr{
				Key:     "$limit",
				Value:   s_limit,
				Message: invalid_limit_err(lang),
			}
		}
	}

	if !ok_page && !ok_limit {
		return 0, 0, false, nil
	}

	if !ok_page {
		return 0, 0, false, &FilterErr{
			Key:     "$page",
			Value:   OmitVal,
			Message: missing_page_num_err(lang),
		}
	}

	if !ok_limit {
		return 0, 0, false, &FilterErr{
			Key:     "$limit",
			Value:   OmitVal,
			Message: missing_limit_err(lang),
		}
	}

	return page, limit, true, nil
}

func invalid_page_num_err(lang string) string {
	if lang == "ar" {
		return "رقم الصفحة غير صالح"
	}
	return "Page number is invalid"
}

func invalid_limit_err(lang string) string {
	if lang == "ar" {
		return "الحد غير صالح"
	}
	return "The limit is invalid"
}

func missing_page_num_err(lang string) string {
	if lang == "ar" {
		return "رقم الصفحة مفقود"
	}
	return "Page number is missing"
}

func missing_limit_err(lang string) string {
	if lang == "ar" {
		return "الحد مفقود"
	}
	return "The limit is missing"
}
