package sequence

import "sync"

//SeqFunc is the type of a function whoes argument is a Sequencable
type SeqFunc func(Sequencable)

//Iterable defines sequence method rules
type Iterable interface {
	Next()
	HasNext() bool
	First() interface{}
	Last() interface{}
	Key() interface{}
	Value() interface{}
	Length() int
	Reset()
}

//Sequencable defines a sequence method rules
type Sequencable interface {
	Iterator() Iterable
	Get(interface{}) interface{}
	Add(...interface{}) Sequencable
	Delete(...interface{}) Sequencable
	Clear() Sequencable
	Obj() interface{}
	Length() int
	Clone() Sequencable
	Mutate(SeqFunc)
	Seq() Sequencable
}

//Sequence is the root level structure for all sequence types
type Sequence struct {
	parent Sequencable
	writer *SeqWriter
}

//SeqWriter represents write operations to be performed on a sequence
//created to avoid race conditions
type SeqWriter struct {
	write chan func()
	lock  *sync.Mutex
}

//Stack adds a function call into the writer stack
func (l *SeqWriter) Stack(fn func()) {
	l.write <- fn
	l.Flush()
}

//Flush begins writing or else ignores if write already started and inprocess
func (l *SeqWriter) Flush() {
	l.lock.Lock()
	for fx := range l.write {
		fx()
	}
	l.lock.Unlock()
}

//ListSequence represents a sequence for arrays,splice type structures
type ListSequence struct {
	*Sequence
	data []interface{}
}

//Seq returns the sequence as a sequencable
func (l *ListSequence) Seq() Sequencable {
	return Sequencable(l)
}

//Get retrieves the value
func (l *ListSequence) Get(d interface{}) interface{} {
	dd, ok := d.(int)

	if !ok {
		return nil
	}

	return l.data[d]
}

//Add for the ListSequence adds all supplied arguments at once to the list
func (l *ListSequence) Add(f ...interface{}) Sequencable {
	l.writer.Stack(func() {
		l.data = append(l.data, f...)
	})
	l.writer.Flush()
	return l.Seq()
}

//Delete for the ListSequence adds all supplied arguments at once to the list
func (l *ListSequence) Delete(f interface{}) Sequencable {
	ind, ok := f.(int)

	if !ok {
		return l.Seq()
	}

	l.writer.Stack(func() {
		l.data = append(l.data[:ind], s.data[ind+1:])
	})

	l.writer.Flush()

	return l.Seq()
}

//ImmutableSequence is the root level of immutable sequence types
type ImmutableSequence struct {
	*Sequence
}

//MapIterator provides an iterator for the map structure
type MapIterator struct {
	Iterable
	data    map[interface{}]interface{}
	updater func(*MapIterator)
}

//GrabKeys returns a list of the given map keys
func GrabKeys(b map[interface{}]interface{}) []interface{} {
	keys := make([]interface{}, len(b))
	count := 0

	for k := range b {
		keys[count] = k
		count++
	}

	return keys
}

//NewMapIterator returns a new mapiterator for use
func NewMapIterator(m map[interface{}]interface{}) *MapIterator {
	keys := GrabKeys(m)
	kit := NewListIterator(keys)

	upd := func(f *MapIterator) {
		keys = GrabKeys(f.data)
		f.Iterable = NewListIterator(keys)
	}

	return &MapIterator{Iterable(kit), m, upd}
}

//NewReverseMapIterator returns a new mapiterator for use
func NewReverseMapIterator(m map[interface{}]interface{}) *MapIterator {
	keys := GrabKeys(m)
	kit := NewReverseListIterator(keys)

	upd := func(f *MapIterator) {
		keys = GrabKeys(f.data)
		f.Iterable = NewReverseListIterator(keys)
	}

	return &MapIterator{Iterable(kit), m, upd}
}

//ListIterator handles interator over arrays,slices
type ListIterator struct {
	data  []interface{}
	index int
}

//First returns the first value of the iterator
func (m *MapIterator) First() interface{} {
	return m.data[m.Iterable.First()]
}

//Last returns the last value of the iterator
func (m *MapIterator) Last() interface{} {
	return m.data[m.Iterable.Last()]
}

//Next moves to the next item
func (m *MapIterator) Next() {
	m.Iterable.Next()
	if m.Iterable.Length() != len(m.data) {
		m.updater(m)
	}
}

//Value returns the current value of the iterator
func (m *MapIterator) Value() interface{} {
	k := m.Key()
	return m.data[k]
}

//Key returns the current key of the iterator
func (m *MapIterator) Key() interface{} {
	return m.Iterable.Value()
}

//ReverseListIterator returns a reverse iterator
type ReverseListIterator struct {
	*ListIterator
}

//NewReverseListIterator returns a new reverse interator
func NewReverseListIterator(b []interface{}) *ReverseListIterator {
	return &ReverseListIterator{NewListIterator(b)}
}

//Key returns the current index of the iterator
func (r *ReverseListIterator) Key() interface{} {
	k, _ := r.ListIterator.Key().(int)

	if k < 0 {
		return nil
	}

	return (len(r.data) - 1) - k
}

//Value returns the value of the data with the index value
func (r *ReverseListIterator) Value() interface{} {
	k, _ := r.Key().(int)

	if k < 0 || k > len(r.data) {
		return nil
	}

	return r.data[k]
}

//NewListIterator returns a new iterator for the []interface{}
func NewListIterator(b []interface{}) *ListIterator {
	return &ListIterator{b, 0}
}

//HasNext calls the next item
func (l *ListIterator) HasNext() bool {
	if l.index <= 0 || l.index < len(l.data) {
		return true
	}

	return false
}

//Next moves to the next item
func (l *ListIterator) Next() {
	if !l.HasNext() {
		return
	}
	l.index++
}

//Length returns the size iterators
func (l *ListIterator) Length() int {
	return len(l.data)
}

//Reset reverst the iterators index
func (l *ListIterator) Reset() {
	l.index = 0
}

//Key returns the current index of the iterator
func (l *ListIterator) Key() interface{} {
	return l.index
}

//Value returns the value of the data with the index value
func (l *ListIterator) Value() interface{} {
	k, _ := l.Key().(int)

	if k < 0 {
		return nil
	}

	return l.data[k]
}

//First returns the current first item
func (l *ListIterator) First() interface{} {
	return l.data[0]
}

//Last returns the last of the data
func (l *ListIterator) Last() interface{} {
	return l.data[len(l.data)-1]
}
