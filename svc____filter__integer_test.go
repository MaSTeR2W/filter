package filter_test

import (
	"net/url"
	"testing"

	"github.com/MaSTeR2W/filter/v1"
)

func TestIntFilter(t *testing.T) {

	var sql = "SELECT * FROM users"

	var fs = filter.NewFilters(
		filter.FilterConfigs{
			SqlSelect: sql,
		},
		filter.NewIntFilter(filter.IntFilterOpts{
			Key:       "age",
			EnableMax: true,
			Max:       10,
			EnableMin: true,
			Min:       2,
		}),
	)

	var langAr = "ar"

	var test_int_filter_eq = func(t *testing.T) {
		var vals = url.Values{
			"age[eq]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE age=8" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_gt = func(t *testing.T) {
		var vals = url.Values{
			"age[gt]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE age>8" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_gte = func(t *testing.T) {
		var vals = url.Values{
			"age[gte]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE age>=8" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_lt = func(t *testing.T) {
		var vals = url.Values{
			"age[lt]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE age<8" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_lte = func(t *testing.T) {
		var vals = url.Values{
			"age[lte]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE age<=8" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_null = func(t *testing.T) {
		var vals = url.Values{
			"age[null]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE age=NULL" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_ltgt = func(t *testing.T) {
		var vals = url.Values{
			"age[lt]": []string{"9"},
			// lt has higher priority than lte
			"age[lte]": []string{"9"},

			"age[gt]": []string{"8"},
			// gt has higher proirity than gte
			"age[gte]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE (age>8 AND age<9)" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_ltegt = func(t *testing.T) {
		var vals = url.Values{
			"age[lte]": []string{"9"},
			"age[gt]":  []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE (age>8 AND age<=9)" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_ltgte = func(t *testing.T) {
		var vals = url.Values{
			"age[lt]": []string{"9"},

			"age[gte]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE (age>=8 AND age<9)" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_ltegte = func(t *testing.T) {
		var vals = url.Values{
			"age[lte]": []string{"9"},
			"age[gte]": []string{"8"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err != nil {
			t.Error("error should be nil: ", err)
			return
		}

		if query != sql+" WHERE (age>=8 AND age<=9)" {
			t.Error("invalid query:", query)
			return
		}
	}

	var test_int_filter_invalid_num_ar = func(t *testing.T) {
		var vals = url.Values{
			"age[eq]": []string{"15.5"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err == nil {
			t.Error("error should not be nil: ", query)
			return
		}

		if err.Error() != "[\nعدد غير صالح\n]" {
			t.Error(err.Error())
			return
		}

		if query != "" {
			t.Error("should return empty query")
			return
		}
	}

	var test_int_filter_large_num_err_ar = func(t *testing.T) {
		var vals = url.Values{
			"age[eq]": []string{"15"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err == nil {
			t.Error("error should not be nil: ", query)
			return
		}

		if err.Error() != "[\nيجب أن يكون العدد أصغر من أو يساوي 10\n]" {
			t.Error(err.Error())
			return
		}

		if query != "" {
			t.Error("should return empty query")
			return
		}
	}

	var test_int_filter_small_num_err_ar = func(t *testing.T) {
		var vals = url.Values{
			"age[eq]": []string{"1"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langAr)

		if err == nil {
			t.Error("error should not be nil: ", query)
			return
		}

		if err.Error() != "[\nيجب أن يكون العدد أكبر من أو يساوي 2\n]" {
			t.Error(err.Error())
			return
		}

		if query != "" {
			t.Error("should return empty query")
			return
		}
	}

	var langEn = "en"

	var test_int_filter_invalid_num_en = func(t *testing.T) {
		var vals = url.Values{
			"age[eq]": []string{"15.5"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langEn)

		if err == nil {
			t.Error("error should not be nil: ", query)
			return
		}

		if err.Error() != "[\ninvalid number\n]" {
			t.Error(err.Error())
			return
		}

		if query != "" {
			t.Error("should return empty query")
			return
		}
	}

	var test_int_filter_large_num_err_en = func(t *testing.T) {
		var vals = url.Values{
			"age[eq]": []string{"18"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langEn)

		if err == nil {
			t.Error("error should not be nil: ", query)
			return
		}

		if err.Error() != "[\nThe number should be less than or equal to 10\n]" {
			t.Error(err.Error())
			return
		}

		if query != "" {
			t.Error("should return empty query")
			return
		}
	}

	var test_int_filter_small_num_err_en = func(t *testing.T) {
		var vals = url.Values{
			"age[eq]": []string{"1"},
		}

		var query, err = fs.ValidateAndConstruct(vals, langEn)

		if err == nil {
			t.Error("error should not be nil: ", query)
			return
		}

		if err.Error() != "[\nThe number should be greater than or equal to 2\n]" {
			t.Error(err.Error())
			return
		}

		if query != "" {
			t.Error("should return empty query")
			return
		}
	}

	t.Run("test_int_filter_eq", test_int_filter_eq)
	t.Run("test_int_filter_gt", test_int_filter_gt)
	t.Run("test_int_filter_gte", test_int_filter_gte)
	t.Run("test_int_filter_lt", test_int_filter_lt)
	t.Run("test_int_filter_lte", test_int_filter_lte)
	t.Run("test_int_filter_null", test_int_filter_null)
	t.Run("test_int_filter_ltgt", test_int_filter_ltgt)
	t.Run("test_int_filter_ltegt", test_int_filter_ltegt)
	t.Run("test_int_filter_ltgte", test_int_filter_ltgte)
	t.Run("test_int_filter_ltegte", test_int_filter_ltegte)
	t.Run("test_int_filter_invalid_num_ar", test_int_filter_invalid_num_ar)
	t.Run("test_int_filter_large_num_err_ar", test_int_filter_large_num_err_ar)
	t.Run("test_int_filter_small_num_err_ar", test_int_filter_small_num_err_ar)
	t.Run("test_int_filter_invalid_num_en", test_int_filter_invalid_num_en)
	t.Run("test_int_filter_large_num_err_en", test_int_filter_large_num_err_en)
	t.Run("test_int_filter_small_num_err_en", test_int_filter_small_num_err_en)
}
