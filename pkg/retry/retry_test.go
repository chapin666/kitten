package retry

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestDoFunc(t *testing.T) {
	count := 0
	err := DoFunc(2, func() error {
		count++
		return errors.New("simulation error")
	}, func(i int) time.Duration {
		if i != 1 {
			t.Error("Number of sleep wrong", i)
		}
		return time.Millisecond
	})

	if count != 2 {
		t.Error("Number of retries wrong", count)
	}

	if err == nil || err != ErrMaxRetries {
		t.Error("exceeded retry limit", err)
	}
}

func ExampleDoFunc() {
	var (
		count int
		value string
	)

	err := DoFunc(3, func() error {
		if count > 1 {
			value = "foo"
			return nil
		}
		count++
		return errors.New("not allowed")
	}, func(i int) time.Duration {
		return time.Millisecond * time.Duration(i)
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(value)
	// Output: foo
}

type tryTest struct {
	count int
}

func (t *tryTest) Try() error {
	t.count = t.count + 1
	return errors.New("simulation error")
}

func (t *tryTest) Count() int {
	return t.count
}

func TestDo(t *testing.T) {
	tt := &tryTest{}
	err := Do(2, tt, func(i int) time.Duration {
		return time.Millisecond * time.Duration(i)
	})

	if v := tt.Count(); v != 2 {
		t.Error("Number of retries wrong", v)
	}

	if err == nil || err != ErrMaxRetries {
		t.Error("exceeded retry limit", err)
	}
}
