package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Valid return true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Mew initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func NewFromMap(data map[string]string) *Form {
	d := url.Values{}

	for key, value := range data {
		d.Add(key, value)
	}

	return &Form{
		d,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	return x != ""
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsFloat64(field string) bool {
	x := f.Get(field)
	_, err := strconv.ParseFloat(x, 64)
	if err != nil {
		f.Errors.Add(field, "This field must be a number")
		return false
	}
	return true
}

func (f *Form) IsDate(field string, layout string) bool {
	x := f.Get(field)
	_, err := time.Parse(layout, x)
	if err != nil {
		f.Errors.Add(field, "This field must be a datetime")
		return false
	}
	return true
}
