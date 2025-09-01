package jsonflex_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/krelinga/go-jsonflex"
)

func TestString(t *testing.T) {
	cases := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name: "Direct Object",
			input: Movie{
				"adult": false,
				"title": nil,
				"genre_ids": jsonflex.Array{
					jsonflex.Number(1),
					jsonflex.Number(2),
					jsonflex.Number(3),
				},
			},
			expected: `{
  Adult: false,
  GenreIDs: [
    0: 1,
    1: 2,
    2: 3,
  ],
  Title: null,
}`,
		},
		{
			name: "Array of Object",
			input: []Movie{
				{
					"adult":     false,
					"title":     nil,
					"genre_ids": jsonflex.Array{jsonflex.Number(1), jsonflex.Number(2), jsonflex.Number(3)},
				},
			},
			expected: `[
  0: {
    Adult: false,
    GenreIDs: [
      0: 1,
      1: 2,
      2: 3,
    ],
    Title: null,
  },
]`,
		},
		{
			name: "Array of Strings",
			input: []string{
				"string1",
				"string2",
				"string3",
			},
			expected: `[
  0: "string1",
  1: "string2",
  2: "string3",
]`,
		},
		{
			name: "Array of int32s",
			input: []int32{
				1,
				2,
				3,
			},
			expected: `[
  0: 1,
  1: 2,
  2: 3,
]`,
		},
		{
			name: "Array of float64s",
			input: []float64{
				1.1,
				2.2,
				3.3,
			},
			expected: `[
  0: 1.1,
  1: 2.2,
  2: 3.3,
]`,
		},
		{
			name:     "Unsupported type",
			input:    int(42),
			expected: "unsupported type: int",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := jsonflex.String(c.input)
			if diff := cmp.Diff(c.expected, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
