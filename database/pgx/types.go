package pgx

type Multi struct {
	arr []interface{}
}

func NewMulti() *Multi {
	return &Multi{}
}

func (m *Multi) Add(item interface{}) *Multi {
	m.arr = append(m.arr, item)
	return m
}

func (m *Multi) Len() int {
	return len(m.arr)
}

func (m *Multi) Arr() []interface{} {
	return m.arr
}

func (m *Multi) Model() interface{} {
	if len(m.arr) > 0 {
		return m.arr[0]
	}
	return nil
}
