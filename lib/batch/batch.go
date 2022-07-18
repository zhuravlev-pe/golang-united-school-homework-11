package batch

import (
	"golang.org/x/sync/errgroup"
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

func getBatch(n int64, pool int64) (res []user) {
	if pool < 1 {
		pool = 1
	}
	if pool > n {
		pool = n
	}
	eg := errgroup.Group{}
	eg.SetLimit(int(pool))
	mx := sync.Mutex{}

	result := make([]user, 0, n)
	for i := int64(0); i < n; i++ {
		i := i
		eg.Go(func() error {
			u := getOne(i)
			mx.Lock()
			result = append(result, u)
			mx.Unlock()
			return nil
		})

	}

	err := eg.Wait()
	if err != nil {
		panic(err)
	}
	return result
}
