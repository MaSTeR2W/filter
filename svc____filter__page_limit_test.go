package filter_test

import (
	"net/url"
	"testing"

	"github.com/MaSTeR2W/filter/v1"
)

func TestPageLimitFilter(t *testing.T) {
	const (
		LANG_AR = "ar"
		LANG_EN = "en"
	)

	var fs = filter.NewFilters(
		filter.FilterConfigs{
			Sql:      "SELECT * FROM users",
			Paginate: true,
			LimitMin: 3,
			LimitMax: 20,
		},
	)

	var test_page_limit = func(t *testing.T) {
		var v = url.Values{
			"$page":  []string{"3"},
			"$limit": []string{"12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error("error should be nil:", err)
			return
		}

		if query != "SELECT * FROM users LIMIT 12 OFFSET 24" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_page_limit", test_page_limit)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_fields = func(
		t *testing.T,
		err *filter.FilterErr,
		key string,
		value any,
		msg string,
	) {
		if err.Key != key {
			t.Error("key should be:", key, "not:", err.Key)
		}

		if err.Value != value {
			t.Error("value should be:", value, "not:", err.Value)
		}

		if err.Message != msg {
			t.Error("message should be:", msg, "not:", err.Message)
		}
	}

	var test_invalid_page = func(t *testing.T) {

		var v = url.Values{
			"$page":  []string{"3.5"},
			"$limit": []string{"12"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$page",
			"3.5",
			"رقم الصفحة غير صالح",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$page",
			"3.5",
			"Page number is invalid",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_invalid_page", test_invalid_page)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_invalid_limit = func(t *testing.T) {

		var v = url.Values{
			"$page":  []string{"3"},
			"$limit": []string{"12.5"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			"12.5",
			"الحد غير صالح",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			"12.5",
			"The limit is invalid",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_invalid_limit", test_invalid_limit)

	//
	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_limit_out_of_range = func(t *testing.T) {

		var v = url.Values{
			"$page":  []string{"3"},
			"$limit": []string{"2"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			2,
			"يجب أن يكون الحد 3 على الأقل",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			2,
			"The limit should be at least 3",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		v = url.Values{
			"$page":  []string{"3"},
			"$limit": []string{"25"},
		}

		query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			25,
			"يجب ألا يتجاوز الحد 20",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			25,
			"The limit should not exceed 20",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_limit_out_of_range", test_limit_out_of_range)

	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_missing_limit = func(t *testing.T) {

		var v = url.Values{
			"$page": []string{"3"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			filter.OmitVal,
			"الحد مفقود",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$limit",
			filter.OmitVal,
			"The limit is missing",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_missing_limit", test_missing_limit)

	//
	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_missing_page = func(t *testing.T) {

		var v = url.Values{
			"$limit": []string{"5"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$page",
			filter.OmitVal,
			"رقم الصفحة مفقود",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should throw error")
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$page",
			filter.OmitVal,
			"Page number is missing",
		)

		if query != "" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_missing_page", test_missing_page)

	//
	//
	//
	//
	//
	//
	//
	//
	//

	var test_missing_page_limit = func(t *testing.T) {

		var v = url.Values{}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error("error should be nil:", err)
			return
		}

		if query != "SELECT * FROM users" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_missing_page_limit", test_missing_page_limit)
}
