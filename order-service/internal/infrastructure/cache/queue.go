package cache

type Queue struct {
	items []interface{}
	front int
	back  int
	size  int
	cap   int
}

func NewQueue(initialCap int) *Queue {
	return &Queue{
		items: make([]interface{}, initialCap),
		front: 0,
		back:  0,
		size:  0,
		cap:   initialCap,
	}
}

func (q *Queue) grow() {
	newCap := q.cap * 2
	newItems := make([]interface{}, newCap)
	for i := 0; i < q.size; i++ {
		newItems[i] = q.items[(q.front+i)%q.cap]
	}
	q.items = newItems
	q.front = 0
	q.back = q.size
	q.cap = newCap
}

func (q *Queue) Push(v interface{}) {
	if q.size == q.cap {
		q.grow()
	}
	q.items[q.back] = v
	q.back = (q.back + 1) % q.cap
	q.size++
}

func (q *Queue) Pop() interface{} {
	if q.size == 0 {
		return nil
	}
	v := q.items[q.front]
	q.front = (q.front + 1) % q.cap
	q.size--
	return v
}

func (q *Queue) Peek() interface{} {
	if q.size == 0 {
		return nil
	}
	return q.items[q.front]
}

func (q *Queue) Len() int {
	return q.size
}
