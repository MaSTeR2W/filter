package filter_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/MaSTeR2W/filter/v1"
)

func TestDateFilter(t *testing.T) {
	const LANG_AR = "ar"
	const LANG_EN = "en"

	var sql = "SELECT * FROM users"
	var fs = filter.NewFilters(
		filter.FilterConfigs{
			SqlSelect: sql,
		},
		filter.MustCreateNewDateFilter(
			filter.DateFilterOpts{
				Key:     "modified",
				After:   "2020-01-02",
				Before:  "2024-05-01",
				NullOpt: true,
			},
		),
	)

	var test_eq = func(t *testing.T) {
		var v = url.Values{
			"modified[eq]": []string{"2024-04-29"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		if query != "SELECT * FROM users WHERE modified='2024-04-29'" {
			t.Error("invalid query:", query)
			return
		}
	}
	t.Run("test_eq", test_eq)

	//
	//
	//
	//
	//
	//

	var test_pr = func(t *testing.T) {
		var v = url.Values{
			"modified[pr]": []string{"2024-04-29"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		if query != "SELECT * FROM users WHERE modified<'2024-04-29'" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_pr", test_pr)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_pre = func(t *testing.T) {
		var v = url.Values{
			"modified[pre]": []string{"2024-04-29"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		if query != "SELECT * FROM users WHERE modified<='2024-04-29'" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_pre", test_pre)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_ps = func(t *testing.T) {
		var v = url.Values{
			"modified[ps]": []string{"2024-04-29"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		if query != "SELECT * FROM users WHERE modified>'2024-04-29'" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_ps", test_ps)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_pse = func(t *testing.T) {
		var v = url.Values{
			"modified[pse]": []string{"2024-04-29"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		if query != "SELECT * FROM users WHERE modified>='2024-04-29'" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_pse", test_pse)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_pse_pre = func(t *testing.T) {
		var v = url.Values{
			"modified[pre]": []string{"2024-04-29"},
			"modified[pse]": []string{"2024-04-25"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		// pre coming first
		if query != "SELECT * FROM users WHERE (modified<='2024-04-29' AND modified>='2024-04-25')" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_pse_pre", test_pse_pre)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_not_null = func(t *testing.T) {
		var v = url.Values{
			"modified[null]": []string{"0"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		if query != "SELECT * FROM users WHERE modified IS NOT NULL" {
			t.Error("invalid query:", query)
			return
		}
	}
	t.Run("test_not_null", test_not_null)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_null = func(t *testing.T) {
		var v = url.Values{
			"modified[null]": []string{"1"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		if query != "SELECT * FROM users WHERE modified IS NULL" {
			t.Error("invalid query:", query)
			return
		}
	}
	t.Run("test_null", test_null)

	//
	//
	//
	//
	//
	//
	//
	//

	var test_pse_not_null = func(t *testing.T) {
		var v = url.Values{
			"modified[pse]":  []string{"2024-04-29"},
			"modified[null]": []string{"0"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		// null does not have any effect
		if query != "SELECT * FROM users WHERE modified>='2024-04-29'" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_pse_not_null", test_pse_not_null)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_pse_pre_null = func(t *testing.T) {
		var v = url.Values{
			"modified[pre]":  []string{"2024-04-29"},
			"modified[pse]":  []string{"2024-04-25"},
			"modified[null]": []string{"1"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error(err)
			return
		}

		// pre coming first
		if query != "SELECT * FROM users WHERE ((modified<='2024-04-29' AND modified>='2024-04-25') OR modified IS NULL)" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_pse_pre_null", test_pse_pre_null)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var checkErrFields = func(
		t *testing.T,
		err *filter.FilterErr,
		key string,
		val string,
		msg string,
		pathFirstEl string,
	) {
		if err.Key != key {
			t.Error("key should be:", key, "not:", err.Key)
		}

		if err.Value != val {
			t.Error("value should be:", val, "not:", err.Value)
		}

		if err.Message != msg {
			t.Error("message should be:", msg, "not:", err.Message)
		}

		if err.Path[0] != pathFirstEl {
			t.Error("path: ", err.Path, ", first element should be:", pathFirstEl)
		}
	}

	var test_invalid_date_ar = func(t *testing.T) {
		var v = url.Values{
			// even if the date part is valid
			// this will return and error
			"modified[eq]": []string{"2087-01-12 15:45:12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			"2087-01-12 15:45:12",
			"التاريخ غير صالح",
			"eq",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_invalid_date_ar", test_invalid_date_ar)

	//
	//
	//
	//
	//
	//
	//
	//

	var test_early_date_ar = func(t *testing.T) {
		var v = url.Values{
			"modified[eq]": []string{"2012-01-12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			"2012-01-12",
			"يجب أن يكون التاريخ بعد (2020-01-02), التاريخ الذي أدخلته (2012-01-12)",
			"eq",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_early_date_ar", test_early_date_ar)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_late_date_ar = func(t *testing.T) {
		var v = url.Values{
			"modified[pr]": []string{"2025-01-12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			"2025-01-12",
			"يجب أن يكون التاريخ قبل (2024-05-01), التاريخ الذي أدخلته (2025-01-12)",
			"pr",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_late_date_ar", test_late_date_ar)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_after_now_date_ar = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: "SELET * FROM users",
			},
			filter.MustCreateNewDateFilter(filter.DateFilterOpts{
				Key:      "modified",
				AfterNow: true, // current date
			}),
		)

		var dateBeforeNow = time.Now().AddDate(-1, -1, -1).Format("2006-01-02")

		var v = url.Values{
			"modified[pr]": []string{dateBeforeNow},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var now = time.Now().Format("2006-01-02")

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			dateBeforeNow,
			"يجب أن يكون التاريخ بعد ("+now+"), التاريخ الذي أدخلته ("+dateBeforeNow+")",
			"pr",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_after_now_date_ar", test_after_now_date_ar)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_before_now_date_ar = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: "SELET * FROM users",
			},
			filter.MustCreateNewDateFilter(filter.DateFilterOpts{
				Key:       "modified",
				BeforeNow: true, // current date
			}),
		)

		var dateAfterNow = time.Now().AddDate(1, 1, 1).Format("2006-01-02")

		var v = url.Values{
			"modified[pr]": []string{dateAfterNow},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var now = time.Now().Format("2006-01-02")

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			dateAfterNow,
			"يجب أن يكون التاريخ قبل ("+now+"), التاريخ الذي أدخلته ("+dateAfterNow+")",
			"pr",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_before_now_date_ar", test_before_now_date_ar)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_invalid_date_en = func(t *testing.T) {
		var v = url.Values{
			// even if the date part is valid
			// this will return and error
			"modified[eq]": []string{"2012-01-12 15:45:12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			"2012-01-12 15:45:12",
			"The date is invalid",
			"eq",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_invalid_date_en", test_invalid_date_en)

	//
	//
	//
	//
	//
	//
	//
	//

	var test_early_date_en = func(t *testing.T) {
		var v = url.Values{
			"modified[eq]": []string{"2012-01-12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			"2012-01-12",
			"The date should be after (2020-01-02), the date you entered (2012-01-12)",
			"eq",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_early_date_en", test_early_date_en)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_late_date_en = func(t *testing.T) {
		var v = url.Values{
			"modified[pr]": []string{"2025-01-12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			"2025-01-12",
			"The date should be before (2024-05-01), the date you entered (2025-01-12)",
			"pr",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_late_date_en", test_late_date_en)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_after_now_date_en = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: "SELET * FROM users",
			},
			filter.MustCreateNewDateFilter(filter.DateFilterOpts{
				Key:      "modified",
				AfterNow: true, // current date
			}),
		)

		var dateBeforeNow = time.Now().AddDate(-1, -1, -1).Format("2006-01-02")

		var v = url.Values{
			"modified[pr]": []string{dateBeforeNow},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var now = time.Now().Format("2006-01-02")

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			dateBeforeNow,
			"The date should be after ("+now+"), the date you entered ("+dateBeforeNow+")",
			"pr",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_after_now_date_en", test_after_now_date_en)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_before_now_date_en = func(t *testing.T) {
		var fs = filter.NewFilters(
			filter.FilterConfigs{
				SqlSelect: "SELET * FROM users",
			},
			filter.MustCreateNewDateFilter(filter.DateFilterOpts{
				Key:       "modified",
				BeforeNow: true, // current date
			}),
		)

		var dateAfterNow = time.Now().AddDate(1, 1, 1).Format("2006-01-02")

		var v = url.Values{
			"modified[pr]": []string{dateAfterNow},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		var now = time.Now().Format("2006-01-02")

		checkErrFields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"modified",
			dateAfterNow,
			"The date should be before ("+now+"), the date you entered ("+dateAfterNow+")",
			"pr",
		)

		if query != "" {
			t.Error("should return empty:", query)

			return
		}
	}
	t.Run("test_before_now_date_en", test_before_now_date_en)

	//
	//
	//
	//
	//
	//
	//
	//
	//
}
