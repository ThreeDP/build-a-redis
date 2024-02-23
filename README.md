[![progress-banner](https://backend.codecrafters.io/progress/redis/13c2921a-683d-485f-8347-71193544c9c6)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)
# GoRedis: A Redis Implementation in Go with Test-Driven Development (TDD)
[![bonus](https://github.com//ThreeDP/build-a-redis/actions/workflows/Tests.yml/badge.svg)](https://github.com//ThreeDP/build-a-redis/actions/workflows/Tests.yml)
[![codecov](https://codecov.io/github/ThreeDP/build-a-redis/graph/badge.svg?token=N9AW6Y3JHP)](https://codecov.io/github/ThreeDP/build-a-redis)

Welcome to GoRedis! This project aims to provide a Redis implementation in Go programming language, following the principles of Test-Driven Development (TDD). Redis is a popular in-memory data store used for caching, session management, and real-time analytics, among other use cases.

## Features

**GoRedis currently supports the following features:**

   - PING: Check if the server is running.
   - ECHO: Return the given string.
   - GET: Retrieve the value of a key. (Work in progress)
   - SET: Set the value of a key. (Work in progress)
   - INFO: Get information and statistics about the server. (Work in progress)
   - Replication: Support for replication. (Work in progress)

## About Test-Driven Development (TDD)

Test-Driven Development is a software development approach where tests are written before the actual code implementation. This ensures that the code is testable, and it helps in defining clear requirements and expectations for the functionality being developed.

## Who to test

**Run Init Tests**
> This command runs the unit tests for the project. Unit tests are designed to test individual components of the codebase in isolation, ensuring that each part works correctly on its own.
```sh
   make unit
```

**Run Unit Test Coverage**
> This command calculates the test coverage for the unit tests. Test coverage indicates the percentage of code that is exercised by the unit tests. It helps identify areas of the codebase that may lack proper testing.
```sh
   make cov
```

**Run Benchmark Tests**
> This command executes benchmark tests for the project. Benchmark tests measure the performance of specific functions or algorithms, helping to identify potential performance bottlenecks or areas for optimization.
```sh
   make bench
```

**Run All Tests Above**
> This command runs all of the tests mentioned above in sequence.
```sh
   make t
```

This is a starting point for Go solutions to the
["Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis).
