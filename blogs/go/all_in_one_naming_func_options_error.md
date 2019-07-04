```go
https://github.com/uber/makisu/blob/master/lib/utils/httputil/httputil.go
```

### Look at the Readable Names !!! Great !!!
```go
type sendOptions struct {
	body          io.Reader
	timeout       time.Duration
	acceptedCodes map[int]bool
	headers       map[string]string
	redirect      func(req *http.Request, via []*http.Request) error
	retry         retryOptions
	transport     http.RoundTripper
	ctx           context.Context

	// This is not a valid http option. It provides a way to override
	// http.Client. This should always used by tests.
	client *http.Client

	// This is not a valid http option. It provides a way to override
	// parts of the url. For example, url.Scheme can be changed from
	// http to https.
	url *url.URL

	// This is not a valid http option. HTTP fallback is added to allow
	// easier migration from http to https.
	// In go1.11 and go1.12, the responses returned when http request is
	// sent to https server are different in the fallback mode:
	// go1.11 returns a network error whereas go1.12 returns BadRequest.
	// This causes TestTLSClientBadAuth to fail because the test checks
	// retry error.
	// This flag is added to allow disabling http fallback in unit tests.
	// NOTE: it does not impact how it runs in production.
	httpFallbackDisabled bool
}

// SendOption allows overriding defaults for the Send function.
type SendOption func(*sendOptions)

// SendNoop returns a no-op option.
func SendNoop() SendOption {
	return func(o *sendOptions) {}
}

// SendBody specifies a body for http request
func SendBody(body io.Reader) SendOption {
	return func(o *sendOptions) { o.body = body }
}

// SendTimeout specifies timeout for http request
func SendTimeout(timeout time.Duration) SendOption {
	return func(o *sendOptions) { o.timeout = timeout }
}

// SendHeaders specifies headers for http request
func SendHeaders(headers map[string]string) SendOption {
	return func(o *sendOptions) { o.headers = headers }
}
```
```go
type retryOptions struct {
	max               int
	interval          time.Duration
	backoffMultiplier float64
	backoffMax        time.Duration
}

// RetryOption allows overriding defaults for the SendRetry option.
type RetryOption func(*retryOptions)

// RetryMax sets the max number of retries.
func RetryMax(max int) RetryOption {
	return func(o *retryOptions) { o.max = max }
}

// RetryInterval sets the interval between retries.
func RetryInterval(interval time.Duration) RetryOption {
	return func(o *retryOptions) { o.interval = interval }
}

// RetryBackoff adds exponential backoff between retries.
func RetryBackoff(backoffMultiplier float64) RetryOption {
	return func(o *retryOptions) { o.backoffMultiplier = backoffMultiplier }
}
```

### Combination **Idea** **Innovate** 
```go
// SendRetry will we retry the request on network / 5XX errors.
func SendRetry(options ...RetryOption) SendOption {
	retry := retryOptions{
		max:               3,
		interval:          250 * time.Millisecond,
		backoffMultiplier: 1, // Defaults with no backoff.
		backoffMax:        30 * time.Second,
	}
	for _, o := range options {
		o(&retry)
	}
	return func(o *sendOptions) { o.retry = retry }
}
```

### Usage - Defaults - Easy for Caller
```go
// Send sends an HTTP request. May return NetworkError or StatusError (see above).
func Send(method, rawurl string, options ...SendOption) (resp *http.Response, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, fmt.Errorf("parse url: %s", err)
	}
	opts := sendOptions{
		body:                 nil,
		timeout:              60 * time.Second,
		acceptedCodes:        map[int]bool{http.StatusOK: true},
		headers:              map[string]string{},
		retry:                retryOptions{max: 1},
		transport:            nil, // Use HTTP default.
		ctx:                  context.Background(),
		url:                  u,
		httpFallbackDisabled: false,
	}
	for _, o := range options {
		o(&opts)
	}

	req, err := newRequest(method, opts)
	if err != nil {
		return nil, err
	}

	client := opts.client
	if client == nil {
		client = &http.Client{
			Timeout:       opts.timeout,
			CheckRedirect: opts.redirect,
			Transport:     opts.transport,
		}
	}

	interval := opts.retry.interval
	for i := 0; i < opts.retry.max; i++ {
		if i > 0 {
			time.Sleep(interval)
			interval = min(
				time.Duration(float64(interval)*opts.retry.backoffMultiplier),
				opts.retry.backoffMax)
		}
		resp, err = client.Do(req)
		// Retry without tls. During migration there would be a time when the
		// component receiving the tls request does not serve https response.
		// TODO (@evelynl): disable retry after tls migration.
		if err != nil && req.URL.Scheme == "https" && !opts.httpFallbackDisabled {
			log.Warnf("Failed to send https request: %s. Retrying with http...", err)
			var httpReq *http.Request
			httpReq, err = newRequest(method, opts)
			if err != nil {
				return nil, err
			}
			httpReq.URL.Scheme = "http"
			resp, err = client.Do(httpReq)
		}
		if err != nil {
			continue
		}
		if resp.StatusCode >= 500 && !opts.acceptedCodes[resp.StatusCode] {
			continue
		}
		break
	}
	if err != nil {
		return nil, NetworkError{err}
	}
	if !opts.acceptedCodes[resp.StatusCode] {
		return nil, NewStatusError(resp)
	}
	return resp, nil
}
```
### More Caller Friendly
```go
// Get sends a GET http request.
func Get(url string, options ...SendOption) (*http.Response, error) {
	return Send("GET", url, options...)
}

// Head sends a HEAD http request.
func Head(url string, options ...SendOption) (*http.Response, error) {
	return Send("HEAD", url, options...)
}

// Post sends a POST http request.
func Post(url string, options ...SendOption) (*http.Response, error) {
	return Send("POST", url, options...)
}

// Put sends a PUT http request.
func Put(url string, options ...SendOption) (*http.Response, error) {
	return Send("PUT", url, options...)
}

// Patch sends a PATCH http request.
func Patch(url string, options ...SendOption) (*http.Response, error) {
	return Send("PATCH", url, options...)
}

// Delete sends a DELETE http request.
func Delete(url string, options ...SendOption) (*http.Response, error) {
	return Send("DELETE", url, options...)
}
```
