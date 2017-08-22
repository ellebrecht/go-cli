package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	context := &Context{}
	context.Args = []*Option{
		&Option{},
	}
	expected := 1
	if context.Count() != expected {
		t.Fatal("expected:", expected, "got:", expected)
	}
}

func TestGetStringList(t *testing.T) {
	valGood := "a,b"
	expected := []string{"a", "b"}
	valMissing := ""
	valBad := 1
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Value: &valGood,
		},
		&Option{
			Name:  "test1",
			Value: &valMissing,
		},
		&Option{
			Name: "test2",
		},
		&Option{
			Name:  "test3",
			Value: &valBad,
		},
		&Option{
			Name:  "test4",
			Value: "bad",
		},
	}
	testCases := []struct {
		index    int
		expected []string
		err      error
	}{
		{0, expected, nil},
		{1, nil, errors.New("test1 is missing")},
		{2, nil, errors.New("test2 is missing")},
		{3, nil, errors.New("test3 value has invalid format")},
		{4, nil, errors.New("test4 should be a pointer")},
	}
	for _, tc := range testCases {
		res, err := context.GetStringList(tc.index)
		if res != nil && tc.expected != nil {
			for i, r := range *res {
				if strings.Compare(r, tc.expected[i]) != 0 {
					t.Fatal("expected:", tc.expected[i], "got:", r)
				}
			}
		} else if err != nil && tc.err != nil {
			if strings.Compare(tc.err.Error(), err.Error()) != 0 {
				t.Fatal("expected:", tc.err.Error(), "got:", err.Error())
			}
		} else {
			t.Fatal("unexpected state")
		}
	}
}

func TestGetStringListForFlag(t *testing.T) {
	expectedVal1 := "expected1"
	expectedVal2 := "expected2"
	expected := []string{expectedVal1, expectedVal2}
	expectedInput := expectedVal1 + "," + expectedVal2
	notExpected := "notExpected1,notExpected2"
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Flag:  "test1",
			Value: &expectedInput,
		},
		&Option{
			Flag:  "test2",
			Value: &notExpected,
		},
	}
	res, err := context.GetStringListForFlag("test1")
	if err != nil {
		t.Fatal("got err:", err)
	}
	for i, val := range *res {
		if strings.Compare(expected[i], val) != 0 {
			t.Fatal("expected:", expected[i], "got:", val)
		}
	}
}

func TestGetDictionary(t *testing.T) {
	valGood := "a,b,c,d"
	valKeyMissing := ",b,c,d"
	valValMissing := "a,,c,d"
	valMismatch := "a,bc,d"
	expected := map[string]string{"a": "b", "c": "d"}
	valMissing := ""
	valBad := 1
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Value: &valGood,
		},
		&Option{
			Name:  "test1",
			Value: &valMissing,
		},
		&Option{
			Name: "test2",
		},
		&Option{
			Name:  "test3",
			Value: &valBad,
		},
		&Option{
			Name:  "test4",
			Value: "bad",
		},
		&Option{
			Name:  "test5",
			Value: &valKeyMissing,
		},
		&Option{
			Name:  "test6",
			Value: &valValMissing,
		},
		&Option{
			Name:  "test7",
			Value: &valMismatch,
		},
	}
	testCases := []struct {
		index    int
		expected map[string]string
		err      error
	}{
		{0, expected, nil},
		{1, nil, errors.New("test1 is missing")},
		{2, nil, errors.New("test2 is missing")},
		{3, nil, errors.New("test3 value has invalid format")},
		{4, nil, errors.New("test4 should be a pointer")},
		{5, nil, errors.New("bad key")},
		{6, nil, errors.New("bad val")},
		{7, nil, errors.New("attribute keys:value count does not match")},
	}
	for _, tc := range testCases {
		res, err := context.GetDictionary(tc.index)
		if res != nil && tc.expected != nil {
			res2 := *res
			for _, r := range res2 {
				if strings.Compare(res2[r], tc.expected[r]) != 0 {
					t.Fatal("expected:", tc.expected[r], "got:", res2[r])
				}
			}
		} else if err != nil && tc.err != nil {
			if strings.Compare(tc.err.Error(), err.Error()) != 0 {
				t.Fatal("expected:", tc.err.Error(), "got:", err.Error())
			}
		} else {
			t.Fatal("unexpected state")
		}
	}
}

func TestGetDictionaryForFlag(t *testing.T) {
	key := "testkey"
	val := "expected"
	expectedInput := key + "," + val
	expected := map[string]string{key: val}
	notExpected := map[string]string{key: "notExpected"}
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Flag:  "test1",
			Value: &expectedInput,
		},
		&Option{
			Flag:  "test2",
			Value: &notExpected,
		},
	}
	res, err := context.GetDictionaryForFlag("test1")
	if err != nil {
		t.Fatal("got err:", err)
	}
	r := *res
	if strings.Compare(expected[key], r[key]) != 0 {
		t.Fatal("expected:", expected[key], "got:", r[key])
	}
}

func TestGetString(t *testing.T) {
	valGood := "test"
	valMissing := ""
	valBad := 1
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Value: &valGood,
		},
		&Option{
			Name:  "test1",
			Value: &valMissing,
		},
		&Option{
			Name: "test2",
		},
		&Option{
			Name:  "test3",
			Value: &valBad,
		},
		&Option{
			Name:  "test4",
			Value: "bad",
		},
	}
	testCases := []struct {
		index    int
		expected interface{}
		err      error
	}{
		{0, &valGood, nil},
		{1, nil, errors.New("test1 is missing")},
		{2, nil, errors.New("test2 is missing")},
		{3, nil, errors.New("test3 value has invalid format")},
		{4, nil, errors.New("test4 should be a pointer")},
	}
	for _, tc := range testCases {
		res, err := context.GetString(tc.index)
		if res != nil && tc.expected != nil {
			if res != tc.expected {
				t.Fatal("expected:", tc.expected, "got:", res)
			}
		} else if err != nil && tc.err != nil {
			if strings.Compare(tc.err.Error(), err.Error()) != 0 {
				t.Fatal("expected:", tc.err.Error(), "got:", err.Error())
			}
		} else {
			t.Fatal("unexpected state")
		}
	}
}

func TestGetStringForFlag(t *testing.T) {
	expected := "testVal1"
	notExpected := "testVal2"
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Flag:  "test1",
			Value: &expected,
		},
		&Option{
			Flag:  "test2",
			Value: &notExpected,
		},
	}
	res, err := context.GetStringForFlag("test1")
	if err != nil {
		t.Fatal("got err:", err)
	}
	if strings.Compare(expected, *res) != 0 {
		t.Fatal("expected:", expected, "got:", *res)
	}
}

func TestGetBool(t *testing.T) {
	valGood := false
	valBad := "test"
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Value: &valGood,
		},
		&Option{
			Name: "test1",
		},
		&Option{
			Name:  "test2",
			Value: &valBad,
		},
		&Option{
			Name:  "test3",
			Value: "bad",
		},
	}
	testCases := []struct {
		index    int
		expected interface{}
		err      error
	}{
		{0, &valGood, nil},
		{1, nil, errors.New("test1 is missing")},
		{2, nil, errors.New("test2 value has invalid format")},
		{3, nil, errors.New("test3 should be a pointer")},
	}
	for _, tc := range testCases {
		res, err := context.GetBool(tc.index)
		if res != nil && tc.expected != nil {
			if res != tc.expected {
				t.Fatal("expected:", tc.expected, "got:", res)
			}
		} else if err != nil && tc.err != nil {
			if strings.Compare(tc.err.Error(), err.Error()) != 0 {
				t.Fatal("expected:", tc.err.Error(), "got:", err.Error())
			}
		} else {
			t.Fatal("unexpected state")
		}
	}
}

func TestGetBoolForFlag(t *testing.T) {
	expected := true
	notExpected := false
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Flag:  "test1",
			Value: &expected,
		},
		&Option{
			Flag:  "test2",
			Value: &notExpected,
		},
	}
	res, err := context.GetBoolForFlag("test1")
	if err != nil {
		t.Fatal("got err:", err)
	}
	if expected != *res {
		t.Fatal("expected:", expected, "got:", *res)
	}
}

func TestGetInt(t *testing.T) {
	valGood := 1
	valBad := "test"
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Value: &valGood,
		},
		&Option{
			Name: "test1",
		},
		&Option{
			Name:  "test2",
			Value: &valBad,
		},
		&Option{
			Name:  "test3",
			Value: "bad",
		},
	}
	testCases := []struct {
		index    int
		expected interface{}
		err      error
	}{
		{0, &valGood, nil},
		{1, nil, errors.New("test1 is missing")},
		{2, nil, errors.New("test2 value has invalid format")},
		{3, nil, errors.New("test3 should be a pointer")},
	}
	for _, tc := range testCases {
		res, err := context.GetInt(tc.index)
		if res != nil && tc.expected != nil {
			if res != tc.expected {
				t.Fatal("expected:", tc.expected, "got:", res)
			}
		} else if err != nil && tc.err != nil {
			if strings.Compare(tc.err.Error(), err.Error()) != 0 {
				t.Fatal("expected:", tc.err.Error(), "got:", err.Error())
			}
		} else {
			t.Fatal("unexpected state")
		}
	}
}

func TestGetIntForFlag(t *testing.T) {
	expected := 1
	notExpected := 0
	context := &Context{}
	context.Args = []*Option{
		&Option{
			Flag:  "test1",
			Value: &expected,
		},
		&Option{
			Flag:  "test2",
			Value: &notExpected,
		},
	}
	res, err := context.GetIntForFlag("test1")
	if err != nil {
		t.Fatal("got err:", err)
	}
	if expected != *res {
		t.Fatal("expected:", expected, "got:", *res)
	}
}
