package filter

import (
	"net/url"
	"slices"
	"strings"
)

type orderer struct {
	cols   []string
	s_cols string
}

type orderer_opts struct {
	cols []string
}

func new_orderer(opts orderer_opts) *orderer {
	return &orderer{
		cols:   opts.cols,
		s_cols: strings.Join(opts.cols, ", "),
	}
}

func (o *orderer) order(v url.Values, lang string) (string, error) {
	var order_by, arrange, ok, err = o.get_order_by_arrange(v, lang)

	if err != nil {
		return "", err
	}

	if !ok {
		return "", nil
	}

	return "ORDER BY " + order_by + " " + arrange, nil
}

func (o *orderer) get_order_by_arrange(v url.Values, lang string) (string, string, bool, error) {
	var order_by string
	var ok bool

	if order_by, ok = get_first_el_if_exists(v, "$order_by"); ok {
		if !slices.Contains(o.cols, order_by) {
			return "", "", false, &FilterErr{
				Key:     "$order_by",
				Value:   order_by,
				Message: is_not_one_of_err(o.s_cols, lang),
			}
		}
	} else {
		return "", "", false, nil
	}

	var arrange string

	if arrange, ok = get_first_el_if_exists(v, "$arrange"); ok {
		if arrange != "DESC" {
			arrange = "ASC"
		}
	} else {
		arrange = "ASC"
	}

	return order_by, arrange, true, nil
}

func is_not_one_of_err(one_of string, lang string) string {
	if lang == "ar" {
		return "يجب اختيار واحد مما يلي: (" + one_of + ")"
	}
	return "Should select one of the following: (" + one_of + ")"
}
