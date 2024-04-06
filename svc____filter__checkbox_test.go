package filter_test

import (
	"net/url"
	"testing"

	"github.com/MaSTeR2W/filter/v1"
)

func TestCheckboxIntFilter(t *testing.T) {
	var sql = "SELECT * FROM users"

	var fs = filter.NewFilters(
		filter.FilterConfigs{
			SqlSelect: sql,
		},

		filter.NewCheckboxIntFilter(
			filter.CheckboxIntFilterOpts{
				Key:  "userStatus",
				Opts: []int{0, 1, 2},
			},
		),
	)

	var test_checkbox_int_filter_in = func(t *testing.T) {
		var v = url.Values{
			"userStatus[in]": []string{"0", "2"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE userStatus IN (0,2)" {
			t.Error(query)
			return
		}

	}

	t.Run("test_checkbox_int_filter_in", test_checkbox_int_filter_in)

	var test_checkbox_int_filter_nin = func(t *testing.T) {
		var v = url.Values{
			"userStatus[nin]": []string{"0", "2"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE userStatus NOT IN (0,2)" {
			t.Error(query)
			return
		}

	}
	t.Run("test_checkbox_int_filter_nin", test_checkbox_int_filter_nin)

	var test_col_alias = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: sql,
			},

			filter.NewCheckboxIntFilter(
				filter.CheckboxIntFilterOpts{
					Key:      "userStatus",
					ColAlias: "alias",
					Opts:     []int{0, 1, 2},
				},
			),
		)

		var v = url.Values{
			"userStatus[in]": []string{"0", "2"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE alias IN (0,2)" {
			t.Error(query)
			return
		}
	}

	t.Run("test_col_alias", test_col_alias)

	var test_null_allowed = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: sql,
			},

			filter.NewCheckboxIntFilter(
				filter.CheckboxIntFilterOpts{
					Key:      "userStatus",
					ColAlias: "alias",
					Opts:     []int{0, 1, 2},
					NullOpt:  true,
				},
			),
		)

		var v = url.Values{
			"userStatus[in]":   []string{"0"},
			"userStatus[null]": []string{"1"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE (alias IN (0) AND alias IS NULL)" {
			t.Error(query)
			return
		}

		v = url.Values{
			"userStatus[nin]":  []string{"0"},
			"userStatus[null]": []string{"0"},
		}
		query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE (alias NOT IN (0) AND alias IS NOT NULL)" {
			t.Error(query)
			return
		}

		v = url.Values{
			"userStatus[null]": []string{"0"},
		}
		query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE alias IS NOT NULL" {
			t.Error(query)
			return
		}
	}

	t.Run("test_null_allowed", test_null_allowed)

	var test_fields = func(
		t *testing.T,
		err *filter.FilterErr,
		key string,
		path int,
		value string,
		msg string,
	) {
		if err.Key != key {
			t.Error(err.Key)
		}

		if err.Path[0] != path {
			t.Error(err.Path)
		}

		if err.Value != value {
			t.Error(err.Value)
		}

		if err.Message != msg {
			t.Error(err.Message)
		}
	}

	const LANG_AR = "ar"

	var test_out_of_opts_err_ar = func(t *testing.T) {
		var v = url.Values{
			"userStatus[in]": []string{"3", "5"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var t_err = (*err.(*filter.FilterErrs))[0].(*filter.FilterErr)

		test_fields(t, t_err, "userStatus", 0, "3", "يجب أن يكون العدد واحد من: (0, 1, 2)")

		if query != "" {
			t.Error("query should be empty not:", query)
		}
	}

	const LANG_EN = "en"

	t.Run("test_out_of_opts_err_ar", test_out_of_opts_err_ar)

	var test_out_of_opts_err_en = func(t *testing.T) {
		var v = url.Values{
			"userStatus[in]": []string{"3", "5"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var t_err = (*err.(*filter.FilterErrs))[0].(*filter.FilterErr)

		test_fields(t, t_err, "userStatus", 0, "3", "The number should be one of: (0, 1, 2)")

		if query != "" {
			t.Error("query should be empty not:", query)
		}
	}

	t.Run("test_out_of_opts_err_en", test_out_of_opts_err_en)

	var test_invalid_opts_err_ar = func(t *testing.T) {
		var v = url.Values{
			"userStatus[in]": []string{"1", "1.5"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var t_err = (*err.(*filter.FilterErrs))[0].(*filter.FilterErr)

		test_fields(t, t_err, "userStatus", 1, "1.5", "هذا ليس عددا")

		if query != "" {
			t.Error("query should be empty not:", query)
		}
	}

	t.Run("test_invalid_opts_err_ar", test_invalid_opts_err_ar)

	var test_invalid_opts_err_en = func(t *testing.T) {
		var v = url.Values{
			"userStatus[in]": []string{"1", "null"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var t_err = (*err.(*filter.FilterErrs))[0].(*filter.FilterErr)

		test_fields(t, t_err, "userStatus", 1, "null", "This is not a number")

		if query != "" {
			t.Error("query should be empty not:", query)
		}
	}

	t.Run("test_invalid_opts_err_en", test_invalid_opts_err_en)
}

func TestCheckboxStrFilter(t *testing.T) {
	var sql = "SELECT * FROM users"

	var fs = filter.NewFilters(
		filter.FilterConfigs{
			SqlSelect: sql,
		},

		filter.NewCheckboxStrFilter(
			filter.CheckboxStrFilterOpts{
				Key:     "column",
				NullOpt: false,
				Opts:    []string{"opt1", "opt2", "opt3"},
			},
		),
	)

	var test_in = func(t *testing.T) {
		var v = url.Values{
			"column[in]": []string{"opt1", "opt3"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE column IN ('opt1','opt3')" {
			t.Error(query)
			return
		}

	}

	t.Run("test_in", test_in)

	var test_nin = func(t *testing.T) {
		var v = url.Values{
			"column[nin]": []string{"opt1", "opt3"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE column NOT IN ('opt1','opt3')" {
			t.Error(query)
			return
		}

	}
	t.Run("test_nin", test_nin)

	var test_col_alias = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: sql,
			},

			filter.NewCheckboxStrFilter(
				filter.CheckboxStrFilterOpts{
					Key:      "column",
					ColAlias: "alias",
					Opts:     []string{"opt'1", "opt2", "opt3"},
					NullOpt:  false,
				},
			),
		)

		var v = url.Values{
			"column[in]": []string{"opt'1", "opt2"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE alias IN ('opt''1','opt2')" {
			t.Error(query)
			return
		}
	}

	t.Run("test_col_alias", test_col_alias)

	var test_null_allowed = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: sql,
			},

			filter.NewCheckboxStrFilter(
				filter.CheckboxStrFilterOpts{
					Key:      "column",
					ColAlias: "alias",
					Opts:     []string{"opt1", "opt2", "opt3"},
					NullOpt:  true,
				},
			),
		)

		var v = url.Values{
			"column[in]":   []string{"opt3"},
			"column[null]": []string{"1"},
		}

		var query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE (alias IN ('opt3') AND alias IS NULL)" {
			t.Error(query)
			return
		}

		v = url.Values{
			"column[nin]":  []string{"opt2"},
			"column[null]": []string{"0"},
		}
		query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE (alias NOT IN ('opt2') AND alias IS NOT NULL)" {
			t.Error(query)
			return
		}

		v = url.Values{
			"column[null]": []string{"0"},
		}
		query, err = fs.ValidateAndConstruct(v, "")

		if err != nil {
			t.Error(err)
			return
		}

		if query != sql+" WHERE alias IS NOT NULL" {
			t.Error(query)
			return
		}
	}

	t.Run("test_null_allowed", test_null_allowed)

	var test_fields = func(
		t *testing.T,
		err *filter.FilterErr,
		key string,
		path int,
		value string,
		msg string,
	) {
		if err.Key != key {
			t.Error(err.Key)
		}

		if err.Path[0] != path {
			t.Error(err.Path)
		}

		if err.Value != value {
			t.Error(err.Value)
		}

		if err.Message != msg {
			t.Error(err.Message)
		}
	}

	const LANG_AR = "ar"

	var test_out_of_opts_err_ar = func(t *testing.T) {
		var v = url.Values{
			"column[in]": []string{"opt1", "opt4"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var t_err = (*err.(*filter.FilterErrs))[0].(*filter.FilterErr)

		test_fields(t, t_err, "column", 1, "opt4", "يجب أن يكون الخيار واحد من: (opt1, opt2, opt3)")

		if query != "" {
			t.Error("query should be empty not:", query)
		}
	}

	const LANG_EN = "en"

	t.Run("test_out_of_opts_err_ar", test_out_of_opts_err_ar)

	var test_out_of_opts_err_en = func(t *testing.T) {
		var v = url.Values{
			"column[in]": []string{"opt4", "opt2"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var t_err = (*err.(*filter.FilterErrs))[0].(*filter.FilterErr)

		test_fields(t, t_err, "column", 0, "opt4", "The option should be one of: (opt1, opt2, opt3)")

		if query != "" {
			t.Error("query should be empty not:", query)
		}
	}

	t.Run("test_out_of_opts_err_en", test_out_of_opts_err_en)

}
