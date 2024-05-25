package options

const (
	// Query one upstream server at a time.
	// AdGuard Home uses its weighted random algorithm to pick the server so that the fastest server is used more often.
	LoadBalance int64 = iota + 1
	// Use parallel queries to speed up resolving by querying all upstream servers simultaneously.
	ParallelRequests
	// Query all DNS servers and return the fastest IP address among all responses.
	// This slows down DNS queries as AdGuard Home has to wait for responses from all DNS servers,
	// but improves the overall connectivity.
	FastestIpAddress
)
