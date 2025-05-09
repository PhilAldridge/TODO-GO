# Todo app in golang

## Introduction
This project was created following the requirements of the main task in the BJSS Go Academy.
It contains several features:
* Three separate modes of storage, in-memory, local JSON file or postgreSQL database.
* The ability to switch freely between these storage methods using the 'mode' flag.
* A V1 REST API performing CRUD operations on the stored Todos.
* A V2 REST API handling multiple users, including user registration and login, as well as user-specific todos.
* bCrypt based password encryption
* JWT based secure authentication of the V2 API using middleware.
* Use of the [context] package to add a TraceID and [slog] to enable traceability of calls through the solution.
* Concurrent reads and concurrent safe write on all storage types. This was originally completed using sync.RWMutex, but has been refactored to user the Actor/Communicating Sequential Processes (CSP) pattern to fit the updated syllabus.
* Parallel tests of both the V1 and V2 APIs to validate that the solution is concurrent safe.
* A separate CLI app that makes calls to either the V1 or V2 API, allowing direct access in the terminal.
* Sensitive data stored in a .env file
* Serves a static page using http.FileServer
* Serves dynamic pages using the Todo data and html/template.
* Graceful shutdown: When the interrupt signal is sent to the application, it stops accepting incoming http requests and finishes resolving all open http requests before shutting down.

## Extension work, not yet implemented:
* No interfaces or receivers: Remove all interfaces or receiver functions from your application and instead use the static-singleton pattern.
* Benchmark: Use benchmark unit tests to determine the performance of your application. Run a separate application to bombard your code with requests and evaluate the performance. 
* pprof: Use the pprof utility to profile application performance. 
* Sharding: Split the todo store "back end" module into a separate executable to the "front end" api/cli/repl/web server, run multiple instances of the back end and distribute traffic from the front end using ring hashing. 
* GRPC: Enable communication between front end and back end with GRPC. 