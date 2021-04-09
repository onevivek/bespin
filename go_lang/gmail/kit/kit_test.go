package kit

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name   string
		origin []string
		want   [][]string
		limit  int
	}{
		{
			name:   "split 6 by 2",
			origin: []string{"a", "b", "c", "d", "e", "f"},
			want:   [][]string{{"a", "b"}, {"c", "d"}, {"e", "f"}},
			limit:  2,
		},
		{
			name:   "split 10 by 3",
			origin: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			want:   [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}, {"10"}},
			limit:  3,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Split(c.origin, c.limit)

			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("want %v, got %v", c.want, got)
			}
		})
	}
}
