package internal

import (
	"reflect"
	"testing"
)

func Test_parseDNSToCheck(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test parsing",
			args: args{
				data: "word1,word2,",
			},
			want: []string{"word1", "word2"},
		},
		{
			name: "test parsing with spaces",
			args: args{
				data: " word1, word2 ",
			},
			want: []string{"word1", "word2"},
		},
		{
			name: "test parsing with empty word",
			args: args{
				data: " word1, word2 ,",
			},
			want: []string{"word1", "word2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDNSToCheck(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDNSToCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args[T comparable] struct {
		haystack []T
		needle   T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "yes",
			args: args[string]{
				haystack: []string{"w1", "w2"},
				needle:   "w1",
			},
			want: true,
		},
		{
			name: "no",
			args: args[string]{
				haystack: []string{"w1", "w2"},
				needle:   "w3",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
