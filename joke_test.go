package jokes

import (
	"net/url"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRequestEmptyQuery(t *testing.T) {
	r := Request{}
	q := url.Values{}
	assert.DeepEqual(t, q, r.Query())
}

func TestRequestQuery(t *testing.T) {
	tests := map[string]struct {
		r Request
		u url.Values
	}{
		"amount": {
			Request{Amount: 5},
			url.Values{"amount": {"5"}},
		},
		"contains": {
			Request{Contains: "some text"},
			url.Values{"contains": {"some text"}},
		},
		"lang": {
			Request{Lang: Pt},
			url.Values{"lang": {"pt"}},
		},
		"safe-mode": {
			Request{Safe: true},
			url.Values{"safe-mode": {""}},
		},
		"type": {
			Request{Type: "twopart"},
			url.Values{"type": {"twopart"}},
		},
		"blacklist": {
			Request{Blacklist: Flags{Nsfw, Political, Explicit}},
			url.Values{"blacklistFlags": {"nsfw,political,explicit"}},
		},
		"id-1": {
			Request{ID: &IDRange{1, 0}},
			url.Values{"idRange": {"1"}},
		},
		"id-2": {
			Request{ID: &IDRange{2, 32}},
			url.Values{"idRange": {"2-32"}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.DeepEqual(t, tt.u, tt.r.Query())
		})
	}
}

func TestRequestURL(t *testing.T) {
	j := Request{
		Category: Categories{Dark, Pun},
	}
	assert.Equal(
		t,
		"https://v2.jokeapi.dev/joke/Dark,Pun",
		j.URL(),
	)
}

func TestLangSetValid(t *testing.T) {
	langs := []Lang{Cs, De, En, Es, Fr, Pt}

	for _, s := range langs {
		t.Run(string(s), func(t *testing.T) {
			var (
				l Lang
				e = l.Set(string(s))
			)
			assert.NilError(t, e)
			assert.Equal(t, l, s)
		})
	}
}

func TestLangSetInvalid(t *testing.T) {
	var (
		l Lang
		e = l.Set("invalid")
	)
	assert.Error(t, e, `invalid lang code: "invalid"`)
	assert.Equal(t, l, Lang(""))
}

func TestFlagsString(t *testing.T) {
	tests := map[string]struct {
		f Flags
		s string
	}{
		"1": {
			Flags{Nsfw},
			"nsfw",
		},
		"multi": {
			Flags{Nsfw, Religious, Political},
			"nsfw,religious,political",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.f.String(), tt.s)
		})
	}
}

func TestFlagsAddValid(t *testing.T) {
	flags := []Flag{
		Nsfw,
		Religious,
		Political,
		Racist,
		Sexist,
		Explicit,
	}

	for _, s := range flags {
		t.Run(string(s), func(t *testing.T) {
			var (
				f = Flags{}
				e = f.Add(string(s))
			)
			assert.NilError(t, e)
			assert.DeepEqual(t, f, Flags{s})
		})
	}
}

func TestFlagsAddInvalid(t *testing.T) {
	var (
		f = Flags{}
		e = f.Add("invalid")
	)
	assert.Error(t, e, `invalid flag: "invalid"`)
	assert.DeepEqual(t, f, Flags{})
}

func TestTypeSetValid(t *testing.T) {
	types := []Type{Single, Twopart}

	for _, s := range types {
		t.Run(string(s), func(t *testing.T) {
			var (
				y Type
				e = y.Set(string(s))
			)
			assert.NilError(t, e)
			assert.Equal(t, y, s)
		})
	}
}

func TestTypeSetInvalid(t *testing.T) {
	var (
		y Type
		e = y.Set("invalid")
	)
	assert.Error(t, e, `invalid type: "invalid"`)
	assert.Equal(t, y, Type(""))
}

func TestIDRangeString(t *testing.T) {
	tests := map[string]struct {
		r IDRange
		s string
	}{
		"0": {
			IDRange{0, 0},
			"0",
		},
		"1": {
			IDRange{42, 0},
			"42",
		},
		"2": {
			IDRange{42, 1024},
			"42-1024",
		},
		"3": {
			*ID(64),
			"64",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.r.String(), tt.s)
		})
	}
}

func TestCategoriesString(t *testing.T) {
	tests := map[string]struct {
		c Categories
		s string
	}{
		"1": {
			Categories{Dark},
			"Dark",
		},
		"multi": {
			Categories{Dark, Programming, Pun},
			"Dark,Programming,Pun",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.c.String(), tt.s)
		})
	}
}

func TestCategoriesAddValid(t *testing.T) {
	cats := []Category{
		Any,
		Christmas,
		Dark,
		Misc,
		Programming,
		Pun,
		Spooky,
	}

	for _, s := range cats {
		t.Run(string(s), func(t *testing.T) {
			var (
				c = Categories{}
				e = c.Add(string(s))
			)
			assert.NilError(t, e)
			assert.DeepEqual(t, c, Categories{s})
		})
	}
}

// TODO: test Response parsing
// TODO: integration test with mock API
// TODO: integration test with real API
