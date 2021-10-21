package jokes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// JokeAPI base URL.
var BaseURL = url.URL{
	Scheme: "https",
	Host:   "v2.jokeapi.dev",
}

// A joke request.
type Request struct {
	Amount    int
	Blacklist Flags
	Category  Categories
	Contains  string
	ID        *IDRange
	Lang      Lang
	Safe      bool
	Type      Type
}

// Convert joke request into a "url.Values" map.
func (j Request) Query() url.Values {
	v := url.Values{}

	if j.Amount > 0 {
		v.Set(KeyAmount, fmt.Sprintf("%d", j.Amount))
	}

	if len(j.Blacklist) > 0 {
		v.Set(KeyBlacklist, j.Blacklist.String())
	}

	if j.Contains != "" {
		v.Set(KeyContains, j.Contains)
	}

	if j.ID != nil {
		v.Set(KeyIDRange, j.ID.String())
	}

	if j.Lang != "" {
		v.Set(KeyLang, string(j.Lang))
	}

	if j.Safe {
		v.Set(KeySafe, "")
	}

	if j.Type != "" {
		v.Set(KeyType, string(j.Type))
	}

	return v
}

// Get the full API URL for the joke. If Category is not set,
// will fetch by default Any categories.
func (j Request) URL() string {
	cat := j.Category

	if len(cat) == 0 {
		cat = Categories{Any}
	}

	url := BaseURL
	url.Path = fmt.Sprintf("joke/%s", cat)
	url.RawQuery = j.Query().Encode()

	return url.String()
}

// Perform the HTTP GET request to fetch the joke by the
// default http client.
func (j Request) Get() ([]Joke, error) {
	return j.GetUsingClient(nil)
}

// Perform the HTTP GET request to fetch the joke by using
// the given http.Client
func (j Request) GetUsingClient(client *http.Client) (r []Joke, e error) {
	var (
		cli = client
		url = j.URL()
		res *http.Response
		jsn []byte
	)

	if cli == nil {
		cli = http.DefaultClient
	}

	if res, e = cli.Get(url); e != nil {
		return
	}

	if jsn, e = ioutil.ReadAll(res.Body); e != nil {
		return
	}

	if r, e = ParseResponse(jsn); e != nil {
		return
	}

	return
}

// Create a new Request struct.
func NewRequest() *Request {
	return &Request{
		Blacklist: Flags{},
		Category:  Categories{},
	}
}

// Represents language returned by JokeAPI.
// See https://jokeapi.dev/#lang
type Lang string

func (s *Lang) Set(l string) error {
	switch Lang(l) {
	case Cs, De, En, Es, Fr, Pt:
	default:
		return fmt.Errorf("invalid lang code: %q", l)
	}

	*s = Lang(l)

	return nil
}

const (
	Cs Lang = "cs" // Czech
	De Lang = "de" // German
	En Lang = "en" // English
	Es Lang = "es" // Spanish
	Fr Lang = "fr" // French
	Pt Lang = "pt" // Portugese
)

// Blacklist flag.
// See https://jokeapi.dev/#blacklist-flags
type Flag string

// Blacklist flags.
type Flags []Flag

// Convert Flags into a space-separated string.
func (s Flags) String() string {
	var (
		b strings.Builder
		c = len(s)
	)

	for i, f := range s {
		b.WriteString(string(f))
		if i < c-1 {
			b.WriteString(",")
		}
	}

	return b.String()
}

// Validate and add given string as a blacklist flag.
// Returns error if the given string is an invalid flag.
func (s *Flags) Add(f string) error {
	switch Flag(f) {
	case Nsfw, Religious, Political, Racist, Sexist, Explicit:
	default:
		return fmt.Errorf("invalid flag: %q", f)
	}

	*s = append((*s), Flag(f))

	return nil
}

const (
	Nsfw      Flag = "nsfw"
	Religious Flag = "religious"
	Political Flag = "political"
	Racist    Flag = "racist"
	Sexist    Flag = "sexist"
	Explicit  Flag = "explicit"
)

// Joke type.
// See https://jokeapi.dev/#joke-type
type Type string

// Set this Type to the given string. Returns error if the
// given string is an invalid joke type.
func (s *Type) Set(t string) error {
	switch Type(t) {
	case Single, Twopart:
	default:
		return fmt.Errorf("invalid type: %q", t)
	}

	*s = Type(t)

	return nil
}

const (
	Single  Type = "single"
	Twopart Type = "twopart"
)

// A joke ID range.
// See https://jokeapi.dev/#idrange-param
type IDRange struct {
	Lower int
	Upper int
}

// Convert ID range into s "number[-number]" string.
func (r IDRange) String() string {
	f := "%[1]d"

	if r.Upper > 0 {
		f = "%[1]d-%[2]d"
	}

	return fmt.Sprintf(f, r.Lower, r.Upper)
}

// Returns an IDRange struct for a single joke ID.
func ID(id int) *IDRange {
	return &IDRange{id, 0}
}

// Represents joke category.
// See https://jokeapi.dev/#categories
type Category string

// Joke categories.
type Categories []Category

// Convert Categories into a space-separated string.
func (s Categories) String() string {
	var (
		b strings.Builder
		c = len(s)
	)

	for i, f := range s {
		b.WriteString(string(f))
		if i < c-1 {
			b.WriteString(",")
		}
	}

	return b.String()
}

// Validate and add given string as a Category. Returns an
// error if given string is an invalid category.
func (s *Categories) Add(c string) error {
	switch Category(c) {
	case Any, Misc, Programming, Dark, Pun, Spooky, Christmas:
	default:
		return fmt.Errorf("invalid category: %q", c)
	}

	*s = append((*s), Category(c))

	return nil
}

const (
	Any         Category = "Any"
	Christmas   Category = "Christmas"
	Dark        Category = "Dark"
	Misc        Category = "Misc"
	Programming Category = "Programming"
	Pun         Category = "Pun"
	Spooky      Category = "Spooky"
)

const (
	KeyAmount    string = "amount"
	KeyBlacklist string = "blacklistFlags"
	KeyContains  string = "contains"
	KeyIDRange   string = "idRange"
	KeyLang      string = "lang"
	KeySafe      string = "safe-mode"
	KeyType      string = "type"
)

// A single joke response from the JokeAPI "joke" endpoint.
type Joke struct {
	Category Category      `json:"category"`
	Delivery string        `json:"delivery"`
	Flags    map[Flag]bool `json:"flags"`
	Id       int           `json:"id"`
	Joke     string        `json:"joke"`
	Lang     Lang          `json:"lang"`
	Safe     bool          `json:"safe"`
	Setup    string        `json:"setup"`
	Type     Type          `json:"type"`
}

func (j Joke) String() string {
	if j.Type == Twopart {
		return fmt.Sprintf("%s\n%s", j.Setup, j.Delivery)
	} else {
		return j.Joke
	}
}

// A multi-joke response if request has Amount > 1.
type Jokes struct {
	Jokes []Joke `json:"jokes"`
}

// An error response.
type ErrorResponse struct {
	Cause    []string `json:"causedBy"`
	Code     int      `json:"code"`
	Info     string   `json:"additionalInfo"`
	Internal bool     `json:"internalError"`
	Message  string   `json:"message"`
	Time     int      `json:"timestamp"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Info)
}

// Parse given response JSON into a slice of Jokes. Returns
// an error if response is an error response.
func ParseResponse(jsn []byte) (j []Joke, e error) {
	var (
		raw   map[string]json.RawMessage
		isErr bool
		joke  Joke
	)

	if e = json.Unmarshal(jsn, &raw); e != nil {
		return
	}

	if err, has := raw["error"]; !has {
		e = errors.New(`response has no "error" property`)
		return
	} else if e = json.Unmarshal(err, &isErr); e != nil {
		return
	}

	if isErr {
		var err ErrorResponse

		if e = json.Unmarshal(jsn, &err); e == nil {
			e = err
		}

		return
	}

	if _, has := raw["amount"]; has {
		var jokes Jokes

		if e = json.Unmarshal(jsn, &jokes); e == nil {
			j = jokes.Jokes
		}

		return
	}

	if e = json.Unmarshal(jsn, &joke); e == nil {
		j = []Joke{joke}
	}

	return
}
