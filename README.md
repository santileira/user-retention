# user-retention


## Overview

This repository runs a script to calculate the user retention based on the application's user activity.

## Requirements

The server is implemented in [Go 1.15](https://go.dev).
To install Go, follow the [instructions](https://go.dev/doc/install).

## Getting Started

- To run the script, go to the root folder and execute:

```
make all
make run
```

`make all` command will test and build the script.

`make run` command will start the script.

If you have docker installed and want to run it in a containerised environment, run:

```
make docker-build
make docker-run
```

- To test the code and see the coverage, go to the root folder and execute:

```
make test
make coverage
```

## Customising the script

There are many arguments to customize the server:

- log-level, by default it's debug. Ex: it can be info.
- file-path, by default it's ./example.csv. If you want to use another file, put it on the root of the repository and 
  pass the file path as argument.

Example:

```
./user-retention script --log-level=info --file-path=./example2.csv 
```

or 

```
docker run -p 8080:8080 -it --rm user-retention --file-path=./example2.csv
```

## Future work

- Improve the algorithm to calculate user retention to be more flexible.
- Calculate more metrics to analyze the business.  
- Improve the tests (add more acceptance and integration tests).
- Improve errors (add more custom errors, be more precise with errors).
- Improve the documentation on the code.
- Improve CI/CD.
