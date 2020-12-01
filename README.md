**Time spent**
Roughly around 5 hours I think, on and off, with maybe an hour wasted on wondering about how many stats to gather for "stats that might be useful for debugging" lol

**Dependencies**
go 1.13

**How to run:**
`go mod tidy`
`go test ./processors`
`go test ./utils`
`go run . -file {the-http-log-file.txt}`

**# How I'd improve it**
Room for improvements -- if I'm spending more time one it -- including: more completed error checking, better naming, and better testing coverage
Though probably off topic -- I'd make it a dockerized server with rest endpoint for file/streaming input and some js with server-send-message for output, more fun imo

Performance wise it's already decent I think ... it's already O(n) + in-place + async + multi-threading. 
And for scaling up -- assuming adding http/rest endpoints on both end (streaming log input and reporting/persisting back-ends )  -- it can just be put on an aws lambda, or get dockerized and run in a k8s deployment