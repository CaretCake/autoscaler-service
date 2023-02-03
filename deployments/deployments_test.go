package autoscaler

import "testing"

func TestAdd(t *testing.T) {
	l, r := 1, 2
	expect := 3

	got := Add(l, r)

	if expect != got {
		t.Errorf("Failed to Add %v and %v. Got %v, expected %v\n",
			l, r, got, expect)
	}
}
