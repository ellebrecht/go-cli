package cli

import "testing"
import "strings"

func TestItemFromRawJSON(t *testing.T) {
	idKey := "id"
	idVal := "test1"
	nameKey := "name"
	nameVal := "test2"
	infoKey := "info"
	infoVal := "test3"
	meta := &Meta{
		RawJSON: `{
			 "` + idKey + `": "` + idVal + `",
			 "` + nameKey + `": "` + nameVal + `",
			 "` + infoKey + `": "` + infoVal + `",
			 "foo": "bar"
		}`,
	}
	item, err := meta.ItemFromRawJSON(idKey, nameKey, infoKey)
	if err != nil {
		t.Fatal("got:", err)
	}
	testCases := []struct {
		val      string
		expected string
	}{
		{item.ID, idVal},
		{item.Name, nameVal},
		{item.Info, infoVal},
	}
	for _, tc := range testCases {
		if strings.Compare(tc.expected, tc.val) != 0 {
			t.Fatal("expected:", tc.expected, "got:", tc.val)
		}
	}
}

func TestItemsFromRawJSON(t *testing.T) {
	idKey := "id"
	idVal := "test1"
	nameKey := "name"
	nameVal := "test2"
	infoKey := "info"
	infoVal := "test3"
	meta := &Meta{
		RawJSON: `[{
			 "` + idKey + `": "` + idVal + `",
			 "` + nameKey + `": "` + nameVal + `",
			 "` + infoKey + `": "` + infoVal + `",
			 "foo": "bar"
		},{
			"` + idKey + `": "` + idVal + `",
			 "` + nameKey + `": "` + nameVal + `",
			 "` + infoKey + `": "` + infoVal + `",
			 "foo": "bar"
		}]`,
	}
	items, err := meta.ItemsFromRawJSON(idKey, nameKey, infoKey)
	if err != nil {
		t.Fatal("got:", err)
	}
	testCases := []struct {
		val      string
		expected string
	}{
		{items[0].ID, idVal},
		{items[0].Name, nameVal},
		{items[0].Info, infoVal},
		{items[1].ID, idVal},
		{items[1].Name, nameVal},
		{items[1].Info, infoVal},
	}
	for _, tc := range testCases {
		if strings.Compare(tc.expected, tc.val) != 0 {
			t.Fatal("expected:", tc.expected, "got:", tc.val)
		}
	}
}

func TestUnmarshalRawJSON(t *testing.T) {
	type data struct {
		Elem string
	}
	expected := "test"
	test := &data{}
	meta := &Meta{
		RawJSON: `{
			 "elem": "` + expected + `"
		}`,
	}
	err := meta.UnmarshalRawJSON(test)
	if err != nil {
		t.Fatal("got:", err)
	}
	if strings.Compare(expected, test.Elem) != 0 {
		t.Fatal("expected:", expected, "got:", test.Elem)
	}
}

func TestUnmarshalRawJSONAtPath(t *testing.T) {
	type data struct {
		Elem string
	}
	expected := "test"
	test := &data{}
	meta := &Meta{
		RawJSON: `{
			 "path1": {
				 "elem": "` + expected + `"
			 }
		}`,
	}
	err := meta.UnmarshalRawJSONAtPath(test, "path1")
	if err != nil {
		t.Fatal("got:", err)
	}
	if strings.Compare(expected, test.Elem) != 0 {
		t.Fatal("expected:", expected, "got:", test.Elem)
	}
}

func TestObjectForPath(t *testing.T) {
	expected := "test"
	meta := &Meta{
		RawJSON: `{
			 "path1": {
				 "path2": {
					 "elem": "` + expected + `"
				 }
			 }
		}`,
	}
	obj, err := meta.ObjectForPath("path1.path2.elem")
	if err != nil {
		t.Fatal("got:", err)
	}
	if strings.Compare(expected, obj.(string)) != 0 {
		t.Fatal("expected:", expected, "got:", obj)
	}
}
