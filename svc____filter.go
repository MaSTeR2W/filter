package filter

import "net/url"

type filters struct {
	sql_select string
	sql_count  string
	filters    []Filter
	len        int
	paginate   bool
	paginator  *paginator
	ordering   bool
	orderer    orderer
}

type FilterConfigs struct {
	SqlSelect string
	SqlCount  string
	Paginate  bool
	LimitMin  int // minimum allowed value for limit
	LimitMax  int // maximum allowed value for limit
	OrderBy   []string
}

func NewFilters(cfg FilterConfigs, fs ...Filter) *filters {

	var f = filters{
		sql_select: cfg.SqlSelect,
		sql_count:  cfg.SqlCount,
		filters:    fs,
		paginate:   cfg.Paginate,
		len:        len(fs),
	}

	if cfg.Paginate {
		var pgOpts = paginator_opts{}
		if cfg.LimitMin < 1 {
			pgOpts.limit_min = 1
		} else {
			pgOpts.limit_min = cfg.LimitMin
		}

		if cfg.LimitMax > 0 {
			pgOpts.enable_limit_max = true
			pgOpts.limit_max = cfg.LimitMax
		}
		f.paginator = new_pagintor(pgOpts)
	}

	if cfg.OrderBy != nil {
		f.ordering = true
		f.orderer = *new_orderer(orderer_opts{
			cols: cfg.OrderBy,
		})
	}

	return &f
}

func (f *filters) ValidateAndConstruct(v url.Values, lang string) (string, error) {

	var errs = make(FilterErrs, 0, f.len)

	var conds = ""
	var cond string

	var err error

	for _, f := range f.filters {
		if cond, err = f.validate_and_construct(v, lang); err != nil {
			errs = append(errs, err)
			continue
		}

		if cond != "" {
			conds += " AND " + cond
		}

	}

	var order_by string

	if f.ordering {
		if order_by, err = f.orderer.order(v, lang); err != nil {
			errs = append(errs, err)
		}
	}

	var limit_offset string

	if f.paginate {
		if limit_offset, err = f.paginator.paginate(v, lang); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return "", &errs
	}

	var query_suffix = ""

	if len(conds) > 0 {
		conds = conds[5:]
		query_suffix = " WHERE " + conds
	}

	if order_by != "" {
		query_suffix += " " + order_by
	}

	if limit_offset != "" {
		query_suffix += " " + limit_offset
	}

	return f.sql_select + query_suffix, nil

}

func (f *filters) ValidateAndConstructWithCount(v url.Values, lang string) (string, string, error) {
	var errs = make(FilterErrs, 0, f.len)

	var conds = ""
	var cond string

	var err error

	for _, f := range f.filters {
		if cond, err = f.validate_and_construct(v, lang); err != nil {
			errs = append(errs, err)
			continue
		}

		if cond != "" {
			conds += " AND " + cond
		}

	}

	var order_by string

	if f.ordering {
		if order_by, err = f.orderer.order(v, lang); err != nil {
			errs = append(errs, err)
		} else {
			order_by = " " + order_by
		}
	}

	var limit_offset string

	if f.paginate {
		if limit_offset, err = f.paginator.paginate(v, lang); err != nil {
			errs = append(errs, err)
		} else {
			limit_offset = " " + limit_offset
		}
	}

	if len(errs) > 0 {
		return "", "", &errs
	}

	var where string

	if len(conds) > 0 {
		conds = conds[5:]
		where = " WHERE " + conds
	}

	return f.sql_select + where + order_by + limit_offset, f.sql_count + where, nil
}
