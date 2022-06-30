package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getUsers(s, n int64, ch chan<- user, wg *sync.WaitGroup) {
	defer wg.Done()
	last := s + n
	for ; s < last; s++ {
		ch <- getOne(s)
	}
}

func getBatch(n int64, pool int64) (res []user) {
	if pool < 1 {
		pool = 1
	}
	if pool > n {
		pool = n
	}
	ch := make(chan user)
	wg := sync.WaitGroup{}
	perPool := n / pool
	var startId int64
	for i := 0; i < int(pool); i++ {
		batchCount := n - startId
		batchCount = min(perPool, batchCount)
		s := startId
		wg.Add(1)
		go getUsers(s, batchCount, ch, &wg)

		startId += perPool
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := make([]user, 0, n)
	for u := range ch {
		result = append(result, u)
	}

	return result
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
