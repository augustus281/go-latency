# Singleflight in Go: A Clean Solution to Cache 

Cache stampede is a common problem in distributed systems, especially when the system is under heavy traffic. It happens when a cached value expires (or is missing), and many users request the same data at the same time. Since the cache no longer has the data, all those requests go directly to the backend service or database. This can overload the backend because it has to handle many identical requests simultaneously.

In Go, the `singleflight` package (from `golang.org/x/sync`) helps solve this problem in a simple way. It makes sure that when multiple requests ask for the same data at the same time, only one actual backend call is made. The other requests wait and receive the same result once that single call finishes. This prevents duplicate work and protects the backend from being overwhelmed.

## Cache Stampede

A cache stampede happens when a cached item expires or is removed, and many users request the same data at the same time. Because the data is no longer in the cache, every request goes to the backend service or database instead. This creates a sudden spike in traffic to the backend. If too many requests arrive at once, it can slow down the system or even cause it to crash.

## A Real-World Scenario

In a real production system, a service running inside a Kubernetes pod used Redis as a distributed cache. The service had an endpoint called `/template-details` that returned template information based on a template ID from each request.

Because templates were updated frequently (about every 30 seconds), caching was very important. It helped reduce the number of direct queries sent to the MySQL database.

During a high-traffic period, a network issue happened between the pod and the Redis server. The service could no longer access Redis, so every request resulted in a cache miss. Instead of reading from the cache, the service had to query the MySQL database for every single request.

This caused a sudden and massive increase in database traffic. Since the database was a single point of failure, it quickly became overloaded and temporarily went down. As a result, the entire system became unavailable.

## The Solution - singleflight

To help prevent this kind of cascading failure, Go provides the `singleflight` package in the `golang.org/x/sync module`. This package makes sure that when multiple requests ask for the same data at the same time, only one real request is sent to the backend. The other requests simply wait and receive the same result once that single request finishes.

This simple but powerful approach can greatly reduce pressure on the backend and help avoid system outages caused by cache stampedes.

However, like any solution, `singleflight` also has trade-offs and things you need to consider before using it.