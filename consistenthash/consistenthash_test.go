package consistenthash

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	m := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))

		return uint32(i)
	})
	//2 4 6 12 14 16 22 24 26
	m.Add("6", "4", "2")

	testCase := map[string]string{
		"2":      "2",
		"13":     "4",
		"114514": "2",
	}

	for key, expect := range testCase {
		if result := m.Get(key); result != expect {
			t.Fatalf("key:%s,expected to match %s,but %s got", key, expect, result)
		}
	}

	// Adds 8, 18, 28
	m.Add("8")

	// 27 should now map to 8.
	testCase["27"] = "8"

	for k, v := range testCase {
		if m.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
