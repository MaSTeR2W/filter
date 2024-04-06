package filter_test

import (
	"net/url"
	"testing"

	"github.com/MaSTeR2W/filter/v1"
)

func TestFilter(t *testing.T) {
	var fs = filter.NewFilters(filter.FilterConfigs{
		SqlSelect: "SELECT * FROM users",
		SqlCount:  "SELECT COUNT(*) AS count FROM users",
		Paginate:  true,
		LimitMin:  3,
		LimitMax:  10,
		OrderBy:   []string{"firstName", "lastName"},
	}, filter.NewStrFilter(filter.StrFilterOpts{
		Key: "firstName",
	}))

	var test_query_count = func(t *testing.T) {
		var v = url.Values{
			"firstName[eq]": []string{"marwan"},
			"$order_by":     []string{"lastName"},
			"$arrange":      []string{"DESC"},
			"$limit":        []string{"10"},
			"$page":         []string{"5"},
		}

		var sel, count, err = fs.ValidateAndConstructWithCount(v, "ar")

		if err != nil {
			t.Error("error should be nil:", err)
			return
		}

		if sel != "SELECT * FROM users WHERE firstName='marwan' ORDER BY lastName DESC LIMIT 10 OFFSET 40" {
			t.Error("invalid sel:", sel)
		}

		if count != "SELECT COUNT(*) AS count FROM users WHERE firstName='marwan'" {
			t.Error("invalid count:", count)
		}
	}

	t.Run("test_query_count", test_query_count)
}
