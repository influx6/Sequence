package sequence

import "testing"

func TestList(t *testing.T) {
	data := []interface{}{1, 32, 56, 7}
	li := NewListIterator(data)

	for ; li.HasNext(); li.Next() {
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

	for ; li.HasNext(); li.Next() {
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

	if data[li.Key()] != li.First() {
		t.Fatal("the first value is incorrect with the map first value", li.Key(), li.First())
	}

	if data[li.Key()] == li.Last() {
		t.Fatal("the last value is incorrect with the map first value", li.Key(), li.First())
	}

	if data[li.Key()] == li.Last() {
		t.Fatal("the last value is incorrect with the map first value", li.Key(), li.First())
	}
}
