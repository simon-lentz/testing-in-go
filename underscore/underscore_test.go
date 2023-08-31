package underscore

import "testing"

func TestCamel(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// Manually entered test cases:
		{"simple", args{"thisIsACamelString"}, "this_is_a_camel_string"},
		{"spaces", args{"string with spaces"}, "string with spaces"},
		{"terminate with capital", args{"endsWithA"}, "ends_with_a"},
	}
	for _, tt := range tests {
		// Running table tests as subtests.
		t.Run(tt.name, func(t *testing.T) {
			if got := Camel(tt.args.str); got != tt.want {
				t.Errorf("Camel() = %v, want %v", got, tt.want)
			}
		})
	}
}
