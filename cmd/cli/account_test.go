package main

import "testing"

func Test_selectAccounts(t *testing.T) {
	selectedIds = map[int]struct{}{}

	result := selectAccounts("1", "7", "asd", " ", "", "0", "-1", "3")

	if len(result) != 3 {
		t.Fatalf("lenght should be 3")
	}

	if result[0] != 1 {
		t.Fatalf("1. result should be 1")
	}

	if result[1] != 3 {
		t.Fatalf("2. result should be 3")
	}

	if result[2] != 7 {
		t.Fatalf("3. result should be 7")
	}
}
