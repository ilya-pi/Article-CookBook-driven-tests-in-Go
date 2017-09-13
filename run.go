package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type DoRetry struct {
	times    int
	interval time.Duration
}

type Transmitter func() (*http.Response, error)
type Validator func([]byte, *http.Response) error

type TestStruct struct {
	desc      string
	transmit  Transmitter
	status    int
	validator Validator
	retry     DoRetry
}

// Override the validator
func (ts TestStruct) Validator(validator Validator) TestStruct {
	ts.validator = validator
	return ts
}

// Override the status
func (ts TestStruct) Status(status int) TestStruct {
	ts.status = status
	return ts
}

func Run(t *testing.T, tests []TestStruct) {
	for _, test := range tests {
		keepTesting := true
		var err error
		for tries := 0; keepTesting; tries++ {
			testName := test.desc
			if tries > 0 {
				time.Sleep(test.retry.interval)
				testName = fmt.Sprintf("%s retry %d of %d", test.desc, tries, test.retry.times)
			}
			t.Run(testName, func(t *testing.T) {
				err = test.run()
				if err == nil {
					keepTesting = false
					return
				}
				if tries >= test.retry.times {
					if cerr, ok := err.(ValidationError); ok {
						if cerr.raw != "" {
							t.Logf("raw: %s", cerr.raw)
						}
						t.Logf("expected: %s", cerr.expected)
						t.Logf("received: %s", cerr.received)
						t.Error(cerr)
					} else {
						t.Error(err)
					}
					keepTesting = false
					return
				}
				// Silently drop err, we don't want retryable tests to fail until the last one
			})
		}
		if t.Failed() {
			t.Fatal()
		}
	}
}

func (test TestStruct) run() error {
	r, err := test.transmit()
	if err != nil {
		return errors.Wrapf(err, "Transmit failed")
	}
	if r == nil {
		return nil
	}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return errors.Wrapf(err, "HTTP Request Body read failed")
	}
	return test.check(body, r)
}

func (test TestStruct) check(body []byte, r *http.Response) error {
	if test.status != 0 && r.StatusCode != test.status {
		return errors.Errorf("Invalid status code: expected %d, got %d", test.status, r.StatusCode)
	}
	if test.validator != nil {
		return test.validator(body, r)
	}
	return nil
}
