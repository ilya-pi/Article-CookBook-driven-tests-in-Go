This is a repo to support a short article I wrote about writing backend table tests.

# CookBook-driven tests in Go

This is a call to test backend implementations utilizing cookbook-driven tests, a remix on top of table tests. If you don't know what table tests are, please first fill in the gap [here](https://blog.golang.org/subtests).

Now consider you have a simple http server that counts numbers. It can accept a new number on an endpoint and return a total sum we have at the moment. Plus two more operations of zeroing the current sum and generating a random number.

Since an example is better then a thousand words, here is how we would test it in the cookbook-driven approach:

```
func TestPositiveNumbers(t *testing.T) {
	dg, cb, cleanup := boot("http://localhost:3000")
	defer cleanup()

	number := dg.PositiveNumber()
	increment := dg.Increment(number)

	Run(t, []TestStruct{
		cb.Ac("Clear all"),
		cb.Add("Add some positive number", number),
		cb.Total("Check the total amount", number),
		cb.Add("Add an incerement", increment),
		cb.Total("Check total amount again", number+increment),
		cb.Ac("Check we can clear after all operations"),
	})
}
```

Here `cb` is `Cookbook` struct, that will hold major operations you can perform with your backend:
```
type CookBook struct {
	addr string
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
```

While `dg`, `DataGenerator` will take care of creating the required primitives in order to communicate with the server:

```
type DataGen struct {
}

func (d *DataGen) PositiveNumber() int {
	return rand.Intn(100)
}

func (d *DataGen) Increment(limit int) int {
	return rand.Intn(limit)
}
```

It is important to note here, that you can have dependencies between different steps of the scenario with the use of state introduced on the test level, as seen here — 

```
func TestRand(t *testing.T) {
	_, cb, cleanup := boot("http://localhost:3000")
	defer cleanup()

	var r int

	Run(t, []TestStruct{
		cb.Rand("Assign a random number", &r),
		cb.Total("Check total is persisted", &r),
	})
}
```

With this you can build as complex scenarios as you want.

Sample repository with working code can be found here — https://github.com/ilya-pi/Article-BackendTableTests

First launch the sample server implementation (a simple calculator):

```
% go run main.go

Calculator is up!
```

And then run tests, your output should look something similar to:

```
% go test 

--- PASS: TestPositiveNumbers (0.00s)
    --- PASS: TestPositiveNumbers/Clear_all (0.00s)
    --- PASS: TestPositiveNumbers/Add_some_positive_number (0.00s)
    --- PASS: TestPositiveNumbers/Check_the_total_amount (0.00s)
    --- PASS: TestPositiveNumbers/Add_an_incerement (0.00s)
    --- PASS: TestPositiveNumbers/Check_total_amount_again (0.00s)
    --- PASS: TestPositiveNumbers/Check_we_can_clear_after_all_operations (0.00s)
```

Note how all the boilerplate or retrying and specifying SLA's on particular endpoints is absolutely abstracted away in `CookBook` and `Run` method.

Our production tests look no more complex in their structure, then those you see here. 

How do your tests look like? Do you like editing them?
