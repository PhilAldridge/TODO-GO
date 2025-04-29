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
* Concurrent reads and concurrent safe write on all storage types.
* Parallel tests of both the V1 and V2 APIs to validate that the solution is concurrent safe.
* A separate CLI app that makes calls to either the V1 or V2 API, allowing direct access in the terminal.
* Sensitive data stored in a .env file
