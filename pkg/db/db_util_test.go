package db

import "testing"

func TestNewSet(t *testing.T) {
	set := NewSet()
	if set.items == nil {
		t.Fatalf("uninitialized map in new set")
	}
	if set.Len() > 0 {
		t.Fatalf("newly created set should have 0 items")
	}
}

func TestSet_Add(t *testing.T) {
	testCases := []struct {
		name        string
		add         int
		expectItems []int
	}{
		{
			name:        "add_first",
			add:         1,
			expectItems: []int{1},
		},
		{
			name:        "add_second",
			add:         2,
			expectItems: []int{1, 2},
		},
		{
			name:        "add_existing",
			add:         2,
			expectItems: []int{1, 2},
		},
		{
			name:        "add_third",
			add:         3,
			expectItems: []int{1, 2, 3},
		},
	}

	set := NewSet()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set.Add(tc.add)
			if set.Len() != len(tc.expectItems) {
				t.Errorf("mismatch in number of items. have %d, expected %d",
					set.Len(), len(tc.expectItems))
			}
			for _, ei := range tc.expectItems {
				if !set.Contains(ei) {
					t.Errorf("item %d expected, but not found in the set", ei)
				}
			}
		})
	}
}
