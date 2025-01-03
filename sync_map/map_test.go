package sync

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"testing"
)

type User struct {
	Name string
}

func TestMap(t *testing.T) {
	m := NewMap[int, *User]()
	var count = 2000000
	var wg errgroup.Group
	wg.SetLimit(100000)
	for i := 0; i < count; i++ {
		var k = i
		wg.Go(func() error {
			m.Set(i, &User{Name: fmt.Sprintf("test_%d", k)})
			fmt.Println(fmt.Sprintf("set %d", k))
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		user := m.Get(i)
		if user == nil {
			fmt.Println(fmt.Sprintf("user is nil %d", i))

		} else {
			fmt.Println(user.Name)
		}

	}
	fmt.Println("ok")

}
