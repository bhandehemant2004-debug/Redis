package util

type Queue [][]string

func (q *Queue) Clear() {
	 *q = nil
}

func NewQueue() Queue {
	return Queue{}
}

func (q *Queue) Push(value []string) {
	*q = append(*q, value)
}

func (q *Queue) Pop() ([]string, bool) {
	if len(*q) == 0 {
		return []string{}, false
	}

	value := (*q)[0]
	*q = (*q)[1:]

	return value, true
}

func (q Queue) IsEmpty() bool {
	return len(q) == 0
}
