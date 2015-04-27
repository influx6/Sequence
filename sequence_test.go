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

// func TestSequence(t *testing.T) {
//
// }
