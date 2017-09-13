package main

import (
	"fmt"
	"net/http"
	"time"
)

type CookBook struct {
	addr string
}

func (c *CookBook) Ac(desc string) TestStruct {
	return TestStruct{
		desc: desc,
		transmit: func() (*http.Response, error) {
			return http.Get(c.addr + "/ac")
		},
		status: 200,
	}
}

func (c *CookBook) Add(desc string, number int) TestStruct {
	return TestStruct{
		desc: desc,
		transmit: func() (*http.Response, error) {
			return http.Get(fmt.Sprintf("%s/add?val=%d", c.addr, number))
		},
		status: 200,
		retry: DoRetry{
			times:    5,
			interval: time.Second,
		},
	}
}

func (c *CookBook) Total(desc string, total *int) TestStruct {
	return TestStruct{
		desc: desc,
		transmit: func() (*http.Response, error) {
			return http.Get(c.addr + "/total")
		},
		validator: TotalValidator(*total),
		status:    200,
	}
}

func (c *CookBook) Rand(desc string, r *int) TestStruct {
	return TestStruct{
		desc: desc,
		transmit: func() (*http.Response, error) {
			return http.Get(c.addr + "/rand")
		},
		validator: RandValidator(r),
		status:    200,
	}
}
