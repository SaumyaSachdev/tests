# Tests to measure latency and throughput for Dqlite
1. raft-latency.go - to measure latency for a given number of requests. Requests are sent sequentially.
2. conc-cores.go - to measure throughput by using multiple threads to send requests to the dqlite cluster. Requests are sent at the same time.
3. raft-concurrent.go - uses a Barrier to send requests at exactly the same time.
