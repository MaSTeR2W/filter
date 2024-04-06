package filter_test

import (
	"net/url"
	"testing"

	"github.com/MaSTeR2W/filter"
)

func TestStringFilterExceedMaxLengthAr(t *testing.T) {
	var lang = "ar"
	var v = url.Values{}
	v.Set("name[eq]", "random")

	const sql = "SELECT * FROM users "
	var f = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: sql,
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key:          "name",
		EnableMaxLen: true,
		MaxLen:       5,
	}))

	var query string
	var err error

	if query, err = f.ValidateAndConstruct(v, lang); err != nil {
		if err.Error() != "[\nيجب تقصير هذا النص إلى 5 من الحروف أو أقل (أنت حاليا تستخدم 6 من الحروف)\n]" {
			t.Error(err)
		}
	}

	if query != "" {
		t.Error("invalid query:", query)
	}
}

func TestStringFilterExceedMaxLengthEn(t *testing.T) {
	var lang = "en"
	var v = url.Values{}
	v.Set("name[eq]", "length7")

	const sql = "SELECT * FROM users "
	var f = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: sql,
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key:          "name",
		EnableMaxLen: true,
		MaxLen:       4,
	}))

	var query string
	var err error

	if query, err = f.ValidateAndConstruct(v, lang); err != nil {
		if err.Error() != "[\nShould shorten this text to 4 characters (you are currently using 7 characters)\n]" {
			t.Error(err)
		}
	}

	if query != "" {
		t.Error("invalid query:", query)
	}
}

func TestStringFilterEQ(t *testing.T) {
	var lang = "ar"
	var v = url.Values{}
	v.Set("name[eq]", "random")

	const sql = "SELECT * FROM users "
	var f = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: sql,
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key:          "name",
		EnableMaxLen: true,
		MaxLen:       10,
	}))

	var query string
	var err error

	if query, err = f.ValidateAndConstruct(v, lang); err != nil {
		t.Error(err)
	}

	if query != sql+" WHERE name='random'" {
		t.Error("invalid query:", query)
	}
}

func TestStringFilterSW(t *testing.T) {
	var lang = "ar"
	var v = url.Values{}
	v.Set("name[sw]", "random")

	const sql = "SELECT * FROM users "
	var f = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: sql,
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key:          "name",
		EnableMaxLen: true,
		MaxLen:       10,
	}))

	var query string
	var err error

	if query, err = f.ValidateAndConstruct(v, lang); err != nil {
		t.Error(err)
	}

	if query != sql+" WHERE name LIKE 'random%'" {
		t.Error("invalid query:", query)
	}
}

func TestStringFilterEW(t *testing.T) {
	var lang = "ar"
	var v = url.Values{}
	v.Set("name[ew]", "random")

	const sql = "SELECT * FROM users "
	var f = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: sql,
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key:          "name",
		EnableMaxLen: true,
		MaxLen:       10,
	}))

	var query string
	var err error

	if query, err = f.ValidateAndConstruct(v, lang); err != nil {
		t.Error(err)
	}

	if query != sql+" WHERE name LIKE '%random'" {
		t.Error("invalid query:", query)
	}
}

func TestStringFilterCT(t *testing.T) {
	var lang = "ar"
	var v = url.Values{}
	v.Set("name[ct]", "random")

	const sql = "SELECT * FROM users "
	var f = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: sql,
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key:          "name",
		EnableMaxLen: true,
		MaxLen:       10,
	}))

	var query string
	var err error

	if query, err = f.ValidateAndConstruct(v, lang); err != nil {
		t.Error(err)
	}

	if query != sql+" WHERE name LIKE '%random%'" {
		t.Error("invalid query:", query)
	}
}

func TestStringFilterNULL(t *testing.T) {
	var lang = "ar"
	var v = url.Values{}
	v.Set("name[null]", "random")

	const sql = "SELECT * FROM users "
	var f = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: sql,
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key:          "name",
		EnableMaxLen: true,
		MaxLen:       10,
	}))

	var query string
	var err error

	if query, err = f.ValidateAndConstruct(v, lang); err != nil {
		t.Error(err)
	}

	if query != sql+" WHERE name=NULL" {
		t.Error("invalid query:", query)
	}
}
