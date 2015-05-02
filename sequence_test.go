package sequence

import "testing"

var (
	data = []interface{}{1, 32, 56, 7}
)

func TestSequence(t *testing.T) {

	incr := NewGenerativeIterator(func(p Iterable) (interface{}, interface{}, error) {
		cur, _ := p.Value().(int)
		key, _ := p.Key().(int)

		if p.Value() == nil {
			return cur, key, nil
		}

		cur++
		key++
		return cur, key, nil
	})

	if incr == nil {
		t.Fatal("Generative is not functioning", incr)
	}

	if incr.Length() != 0 {
		t.Fatal("Generative Length is above 0 and already used", incr.Length(), incr)
	}

	pre := 0

	for incr.HasNext() {
		pv, _ := incr.Value().(int)
		if pv >= 10 {
			break
		}

		err := incr.Next()

		if incr.Value() != pre {
			t.Fatal("Incrementing value is not accurate:", err, pre, incr.Value(), incr.Key())
		}

		pre++
	}
}

func TestEmptyList(t *testing.T) {
	li := NewListIterator(make([]interface{}, 0))

	//will not work
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

func TestList(t *testing.T) {
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

func TestBaseIterator(t *testing.T) {
	li := NewListIterator(data)
	bl := IdentityIterator(li)

	if bl == nil {
		t.Fatal("BaseIterator can not be equal to its source", li, bl)
	}

	for bl.HasNext() {
		err := bl.Next()

		if err != nil {
			t.Fatal("Error occcured with reverse list", err)
			break
		}

		ind, _ := bl.Key().(int)
		if bl.Value() != data[ind] {
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
	ls := NewListSequence(nil, 3)

	if ls == nil {
		t.Fatal("recevied nil instead of ListSequence", ls)
	}

	if ls.Length() != 0 {
		t.Fatal("length of new list sequence is not 0", ls.Length(), ls)
	}

	ls.Add(1, 2, 4, 5)

	incrementd := func(n int) {
		i := n
		size := n * 2
		for i < size {
			ls.Add(i)
			i++
		}
	}

	go incrementd(1)
	go incrementd(3)

	if ls.Length() != 4 {
		t.Fatal("even after adding 4 items, list is at 4 size", ls.Length(), ls)
	}

	// using the index
	ls.Delete(2)

	if ls.Length() != 3 {
		t.Fatal("even after deleting 4 items, list is empty", ls.Length(), ls)
	}

	t2 := ls.Get(2)

	if t2 != 5 {
		t.Fatalf("after removal ,value at index %d should be 5 but it is %d", 2, t2)
	}

	cl := ls.Clone()

	if cl.Length() != ls.Length() {
		t.Fatalf("clone must be the same length as origin")
	}

	if cl.Get(0) != ls.Get(0) {
		t.Fatalf("clone first index is not equal with source")
	}

	if ls.RootSeq() != ls {
		t.Fatal("rootseq() must returns itself")
	}

	if ls.Seq() != ls {
		t.Fatal("Seq() must wrap itself and must be equal itself")
	}
}

func TestMapSequence(t *testing.T) {
	ls := NewMapSequence(nil, 0)

	if ls == nil {
		t.Fatal("recevied nil instead of ListSequence", ls)
	}

	if ls.Length() != 0 {
		t.Fatal("length of new list sequence is not 0", ls.Length(), ls)
	}

	ls.Add(1, 'a')
	ls.Add(2, 'b')
	ls.Add(3, 'c')

	if ls.Length() != 3 {
		t.Fatal("even after adding 4 items, list is empty", ls.Length(), ls)
	}

	ls.Delete(2)

	if ls.Length() != 2 {
		t.Fatal("even after deleting 4 items, list is empty", ls.Length(), ls)
	}

	t2 := ls.Get(2)

	if t2 != nil {
		t.Fatalf("after removal ,value at index %d should be 5 but it is %d", 2, t2)
	}

	cl := ls.Clone()

	if cl.Length() != ls.Length() {
		t.Fatalf("clone must be the same length as origin")
	}

	if cl.Get(1) != ls.Get(1) {
		t.Fatalf("clone first index is not equal with source")
	}

	if ls.RootSeq() != ls {
		t.Fatal("rootseq() must returns itself")
	}

	if ls.Seq() != ls {
		t.Fatal("Seq() must wrap itself and must be equal itself")
	}
}
