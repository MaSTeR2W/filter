package filter_test

import (
	"net/url"
	"testing"

	"github.com/MaSTeR2W/filter"
)

func TestOrderBy(t *testing.T) {
	const (
		LANG_AR = "ar"
		LANG_EN = "en"
	)

	var fs = filter.NewFilters(
		filter.FilterConfigs{
			SqlSelect: "SELECT * FROM users",
			OrderBy:   []string{"firstName", "lastName"},
		},
	)

	var test_order_by = func(t *testing.T) {
		var v = url.Values{
			"$order_by": []string{"lastName"},
			// event if arrange is missing or invalid
			// it will take the default value (ASC)
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error("error should be nil:", err)
			return
		}

		if query != "SELECT * FROM users ORDER BY lastName ASC" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_order_by", test_order_by)

	//
	//
	//
	//
	//
	//

	var test_order_by_empty = func(t *testing.T) {
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

	t.Run("test_order_by_empty", test_order_by_empty)

	//
	//
	//
	//
	//
	//

	var test_order_by_asc_desc = func(t *testing.T) {
		var v = url.Values{
			"$order_by": []string{"firstName"},
			"$arrange":  []string{"ASC"},
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error("error should be nil:", err)
			return
		}

		if query != "SELECT * FROM users ORDER BY firstName ASC" {
			t.Error("invalid query:", query)
			return
		}

		v = url.Values{
			"$order_by": []string{"firstName"},
			"$arrange":  []string{"DESC"},
		}

		query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err != nil {
			t.Error("error should be nil:", err)
			return
		}

		if query != "SELECT * FROM users ORDER BY firstName DESC" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_order_by_asc_desc", test_order_by_asc_desc)

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

	//
	//
	//
	//
	//
	//

	var test_order_by_none_exist_col = func(t *testing.T) {
		var v = url.Values{
			"$order_by": []string{"notExist"},
			// event if arrange is missing or invalid
			// it will take the default value (ASC)
		}

		var query, err = fs.ValidateAndConstruct(v, LANG_AR)

		if err == nil {
			t.Error("should return error:", err)
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$order_by",
			"notExist",
			"يجب اختيار واحد مما يلي: (firstName, lastName)",
		)
		if query != "" {
			t.Error("invalid query:", query)
			return
		}

		v = url.Values{
			"$order_by": []string{"notExist"},
			// event if arrange is missing or invalid
			// it will take the default value (ASC)
		}

		query, err = fs.ValidateAndConstruct(v, LANG_EN)

		if err == nil {
			t.Error("should return error:", err)
			return
		}

		test_fields(
			t,
			(*err.(*filter.FilterErrs))[0].(*filter.FilterErr),
			"$order_by",
			"notExist",
			"Should select one of the following: (firstName, lastName)",
		)
		if query != "" {
			t.Error("invalid query:", query)
			return
		}
	}

	t.Run("test_order_by_none_exist_col", test_order_by_none_exist_col)

	//
	//
	//
	//
	//
	//
}
