package sequence

import "testing"

func TestList(t *testing.T) {
	data := []interface{}{1, 32, 56, 7}
	li := NewListIterator(data)

	for li.HasNext() {
		err := li.Next()

		if err != nil {
			t.Fatal("Error occcured with reverse list", err)
			break
		}

		ind, _ := li.Key().(int)
		if li.Value() != data[ind] {
			t.Fatal("Index and value incorrect with list", li.Key(), li.Value(), data)
			break
		}
	}
}

func TestReverseList(t *testing.T) {
	data := []interface{}{1, 32, 56, 7}
	li := NewReverseListIterator(data)

	for li.HasNext() {
		err := li.Next()

		if err != nil {
			t.Fatal("Error occcured with reverse list", err)
			break
		}

		ind, _ := li.Key().(int)
		if li.Value() != data[ind] {
			t.Fatal("Index and value incorrect with list", li.Key(), li.Value(), data)
			break
		}
	}

}

func TestMap(t *testing.T) {
	data := map[interface{}]interface{}{1: "a", 32: "v", 56: "h"}
	li := NewMapIterator(data)

	err := li.Next()

	if data[li.Key()] != li.Value() {
		t.Fatal("the first value is incorrect with the map first value", li.Key(), li.Value())
	}

	err = li.Next()

	if data[li.Key()] != li.Value() {
		t.Fatal("the first value is incorrect with the map first value", li.Key(), li.Value())
	}

	err = li.Next()

	if data[li.Key()] != li.Value() {
		t.Fatal("the first value is incorrect with the map first value", li.Key(), li.Value())
	}

	err = li.Next()

	if err != ErrENDINDEX {
		t.Fatal("error is incorrect", err, "expecting ", ErrENDINDEX)
	}

}

func TestListSequence(t *testing.T) {
	ls := NewListSequence(nil, 0)

	if ls == nil {
		t.Fatal("recevied nil instead of ListSequence", ls)
	}

	if ls.Length() != 0 {
		t.Fatal("length of new list sequence is not 0", ls.Length(), ls)
	}

	ls.Add(1, 2, 4, 5)

	if ls.Length() != 4 {
		t.Fatal("even after adding 4 items, list is empty", ls.Length(), ls)
	}

	ls.Delete(2)

	if ls.Length() != 3 {
		t.Fatal("even after deleting 4 items, list is empty", ls.Length(), ls)
	}

	t2 := ls.Get(2)

	if t2 != 5 {
		t.Fatalf("after removal ,value at index %d should be 5 but it is %d", 2, t2)

	}
}
