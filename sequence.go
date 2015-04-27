package sequence

import (
	"errors"
	"sync"
)

const (
	//MINBUFF states the default minimum buffer size for the write channels
	MINBUFF = 20
)

var (
	//ErrBADValue represents a bad value calculation by the iterator
	ErrBADValue = errors.New("Bad Value!")
	//ErrBADINDEX represents a bad index counter by the iterator
	ErrBADINDEX = errors.New("BadIndex!")
	//ErrENDINDEX represents a reaching of the end of an iterator
	ErrENDINDEX = errors.New("EndIndex!")
)

//MutFunc is the type of a function whoes argument is a Sequencable
type MutFunc func(f interface{}) interface{}

//ProcFunc is the type of a function giving to a BaseIterator
type ProcFunc func(f Iterable) (interface{}, interface{}, error)

//Iterable defines sequence method rules
type Iterable interface {
	Next() error
	HasNext() bool
	Key() interface{}
	Value() interface{}
	Reset()
	Length() int
}

//Sequencable defines a sequence method rules
type Sequencable interface {
	Iterator() Iterable
	Get(interface{}) interface{}
	Clear() Sequencable
	Length() int
	Mutate(MutFunc)
	Clone() Sequencable
	Seq() Sequencable
	Add(...interface{}) Sequencable
	Delete(...interface{}) Sequencable
}

//ListSequencable defines ListSequence method rules
type ListSequencable interface {
	// Sequencable
	Obj() []interface{}
}

//ImmutableSequence is the root level of immutable sequence types
type ImmutableSequence struct {
	*Sequence
}

//Sequence is the root level structure for all sequence types
type Sequence struct {
	parent Sequencable
	writer *SeqWriter
}

//Iterator returns the iterator of the sequence
func (s *Sequence) Iterator() Iterable {
	return s.parent.Iterator()
}

//Length returns the length of the sequence
func (s *Sequence) Length() int {
	return s.parent.Length()
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

//NewSeqWriter returns a new Sequence writer for concurrent use
func NewSeqWriter(size int) *SeqWriter {
	return &SeqWriter{
		make(chan func(), size),
		new(sync.Mutex),
	}
}

//NewBaseSequence returns a base sequence struct
func NewBaseSequence(buff int, parent Sequencable) *Sequence {
	if buff < MINBUFF {
		buff = MINBUFF
	}

	return &Sequence{
		parent,
		NewSeqWriter(buff),
	}
}

//NewListSequence returns a new ListSequence
func NewListSequence(data []interface{}, buff int) *ListSequence {
	if data == nil {
		data = make([]interface{}, 0)
	}

	return &ListSequence{
		NewBaseSequence(buff, nil),
		data,
		buff,
	}
}

//NewMapSequence returns a new MapSequence
func NewMapSequence(data map[interface{}]interface{}, buff int) *MapSequence {
	if data == nil {
		data = make(map[interface{}]interface{})
	}

	return &MapSequence{
		NewBaseSequence(buff, nil),
		data,
		buff,
	}
}

//MapSequence represents a sequence for maps
type MapSequence struct {
	*Sequence
	data   map[interface{}]interface{}
	buffer int
}

//Mutate allows mutation on sequence data
func (l *MapSequence) Mutate(fn MutFunc) {
	l.writer.Stack(func() {
		res, ok := fn(l.data).(map[interface{}]interface{})

		if !ok {
			return
		}

		l.data = res
	})
	l.writer.Flush()
}

//Iterator returns the sequence data iterator
func (l *MapSequence) Iterator() Iterable {
	return NewMapIterator(l.data)
}

//Seq returns the sequence as a sequencable
func (l *MapSequence) Seq() Sequencable {
	return Sequencable(l)
}

//Get retrieves the value
func (l *MapSequence) Get(d interface{}) interface{} {
	return l.data[d]
}

//Clone copies internal structure data
func (l *MapSequence) Clone() Sequencable {
	// l.data = make([]interface{}, 0)
	nd := make(map[interface{}]interface{})

	for k, v := range l.data {
		nd[k] = v
	}

	return NewMapSequence(nd, l.buffer)
}

//Clear wipes internal structure data
func (l *MapSequence) Clear() Sequencable {
	l.data = make(map[interface{}]interface{})
	return l.Seq()
}

//Length returns length of data
func (l *MapSequence) Length() int {
	return len(l.data)
}

//Add for the ListSequence adds all supplied arguments at once to the list
func (l *MapSequence) Add(f ...interface{}) Sequencable {
	l.writer.Stack(func() {
		key := f[0]
		val := f[1]
		l.data[key] = val
	})
	l.writer.Flush()
	return l.Seq()
}

//Delete for the ListSequence adds all supplied arguments at once to the list
func (l *MapSequence) Delete(f ...interface{}) Sequencable {
	for _, v := range f {
		l.writer.Stack(func() {
			_, ok := l.data[v]

			if !ok {
				return
			}

			delete(l.data, v)
		})
	}

	l.writer.Flush()

	return l.Seq()
}

//ListSequence represents a sequence for arrays,splice type structures
type ListSequence struct {
	*Sequence
	data   []interface{}
	buffer int
}

//Mutate allows mutation on sequence data
func (l *ListSequence) Mutate(fn MutFunc) {
	l.writer.Stack(func() {
		res, ok := fn(l.data).([]interface{})

		if !ok {
			return
		}

		l.data = res
	})
	l.writer.Flush()

}

//Iterator returns the sequence data iterator
func (l *ListSequence) Iterator() Iterable {
	return NewListIterator(l.data)
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

	return l.data[dd]
}

//Clone copies internal structure data
func (l *ListSequence) Clone() Sequencable {
	// l.data = make([]interface{}, 0)
	nd := make([]interface{}, l.Length())
	copy(nd, l.data)
	return NewListSequence(nd, l.buffer)
}

//Clear wipes internal structure data
func (l *ListSequence) Clear() Sequencable {
	l.data = make([]interface{}, 0)
	return l.Seq()
}

//Length returns length of data
func (l *ListSequence) Length() int {
	return len(l.data)
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
func (l *ListSequence) Delete(f ...interface{}) Sequencable {
	for _, v := range f {
		ind, ok := v.(int)

		if !ok {
			return l.Seq()
		}

		l.writer.Stack(func() {
			l.data = append(l.data[:ind], l.data[ind+1:])
		})

	}
	l.writer.Flush()

	return l.Seq()
}

//MapIterator provides an iterator for the map structure
type MapIterator struct {
	Iterable
	data    map[interface{}]interface{}
	updater func(*MapIterator) int
	size    int
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

	upd := func(f *MapIterator) int {
		keys = GrabKeys(f.data)
		f.Iterable = NewListIterator(keys)
		return len(keys)
	}

	return &MapIterator{Iterable(kit), m, upd, 0}
}

//NewReverseMapIterator returns a new mapiterator for use
func NewReverseMapIterator(m map[interface{}]interface{}) *MapIterator {
	keys := GrabKeys(m)
	kit := NewReverseListIterator(keys)

	upd := func(f *MapIterator) int {
		keys = GrabKeys(f.data)
		f.Iterable = NewReverseListIterator(keys)
		return len(keys)
	}

	return &MapIterator{Iterable(kit), m, upd, 0}
}

//BaseIterator handles interation over an iterator
type BaseIterator struct {
	parent Iterable
	value  interface{}
	index  interface{}
	proc   ProcFunc
}

//IdentityIterator takes an Iterable and returns an iterator that simple returns
//the root iterators key and value without change,useful for IteratorSequence
func IdentityIterator(b Iterable) *BaseIterator {
	return NewBaseIterator(b, func(root Iterable) (interface{}, interface{}, error) {
		return root.Value(), root.Key(), nil
	})
}

//NewBaseIterator returns a base iterator based on a function evaluator
func NewBaseIterator(b Iterable, fn ProcFunc) *BaseIterator {
	return &BaseIterator{
		b,
		nil,
		nil,
		fn,
	}
}

//HasNext calls the next item
func (l *BaseIterator) HasNext() bool {
	return l.parent.HasNext()
}

//Next moves to the next item
func (l *BaseIterator) Next() error {
	err := l.parent.Next()

	if err == ErrBADValue {
		l.value = nil
		l.index = nil
		return ErrBADValue
	}

	if err != nil {
		return err
	}

	v, k, err := l.proc(l.parent)

	if err != nil {
		return err
	}

	l.value = v
	l.index = k
	return nil
}

//Reset reverst the iterators index
func (l *BaseIterator) Reset() {
	l.parent.Reset()
	l.value = nil
	l.index = nil
}

//Key returns the current index of the iterator
func (l *BaseIterator) Key() interface{} {
	return l.index
}

//Value returns the value of the data with the index value
func (l *BaseIterator) Value() interface{} {
	return l.value
}

//Length returns the parent iterators targets length,not its operation length
func (l *BaseIterator) Length() int {
	return l.parent.Length()
}

//ListIterator handles interator over arrays,slices
type ListIterator struct {
	data  []interface{}
	index int
}

//Next moves to the next item
func (m *MapIterator) Next() error {
	err := m.Iterable.Next()
	if m.size != len(m.data) {
		m.size = m.updater(m)
	}
	return err
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

//Length returns the iterators targets length,not its operation length
func (m *MapIterator) Length() int {
	return len(m.data)
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
	if l.index <= 0 || l.index < (len(l.data)-1) {
		return true
	}
	return false
}

//Next moves to the next item
func (l *ListIterator) Next() error {
	if !l.HasNext() {
		return ErrENDINDEX
	}
	l.index++
	return nil
}

//Length returns the iterators targets length,not its operation length
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
