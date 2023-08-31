package underscore

import "testing"

func TestCamel_Manual(t *testing.T) {
	testCases := []struct {
		arg  string
		want string
	}{
		{"thisIsACamelString", "this_is_a_camel_string"},
		{"string with spaces", "string with spaces"},
		{"endsWithA", "ends_with_a"},
	}
	for _, tt := range testCases {
		t.Logf("Testing: Camel(%q)", tt.arg)
		got := Camel(tt.arg)
		if got != tt.want {
			t.Errorf("Camel(%q) = %q; want %q", tt.arg, got, tt.want)
		}
	}
}
