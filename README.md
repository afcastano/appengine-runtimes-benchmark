# appengine-runtimes-benchmark
Simple benchmark of the different appengine runtimes and frameworks.

Summary:
- Appengine provides several runtimes.
- Each language has pros and cons.
- This is only looking at performance.
- Goal is to help decide the performance tradeoff when choosing a runtime.
- Includes OSS libraries that we use in GO projects as well as vanilla runtimes.

Metodology:
- Each runtime was loaded with 4, 40 and 400 requests per second. Each request will run a Query which will filter and limit to 10 entities.
- This will run for about 10 to 12 mins.
- The data is collected on the appengine dashboard and the results are taken halfway through the test.
- The runtime/library is deployed with the basic stack for each one of them.
- Instance type is F2 for each of them.
- US/Central region

Results:
The results are as follows:

Analysis:
Low load 4 req/sec
====================
- Instance are all the same
- Startup time is key. Nest - Spring are slow. Go is too fast.
- In lower request per second, might be important considering the memory that is required to run the app, since it can trigger another instance.
