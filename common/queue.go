package common

import (
	"container/list"
	"sync"
)

// 先进先出队列
type Queue struct {
	data *list.List
	mut  *sync.RWMutex
}

// 工厂函数 生成 `先进先出队列`
func NewQueue() *Queue {
	return &Queue{data: list.New(), mut: new(sync.RWMutex)}
}

// 入队操作
func (q *Queue) Put(v interface{}) {
	defer q.mut.Unlock()
	q.mut.Lock()
	q.data.PushFront(v)
}

// 出队操作
func (q *Queue) Get() (interface{}, bool) {
	defer q.mut.Unlock()
	q.mut.Lock()
	if q.data.Len() > 0 {
		iter := q.data.Back()
		v := iter.Value
		q.data.Remove(iter)
		return v, true
	}
	return nil, false
}

// 返回队列长度
func (q *Queue) Qsize() int {
	defer q.mut.RUnlock()
	q.mut.RLock()
	return q.data.Len()
}

func (q *Queue) Clear() {
	defer q.mut.Unlock()
	q.mut.Lock()
	for ;q.data.Len()>0;{
		iter := q.data.Back()
		q.data.Remove(iter)
	}
}

// 判断队列是否为空
func (q *Queue) IsEmpty() bool {
	defer q.mut.RUnlock()
	q.mut.RLock()
	return !(q.data.Len() > 0)
}
