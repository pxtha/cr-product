# CR-Product

CR-Product is a Go project with a focus on product management.

## Overall Architecture

![system](system.svg)


## Project Structure


The project is structured into several directories:

- `cmd/server`: Contains the main application entry point.
- `conf`: Contains configuration files.
- `deploy`: Contains Docker and Docker Compose files for deployment.
- `internal/app`: Contains the application logic, including models, routes, and workers.
- `internal/pkg`: Contains packages used across the application.
- `internal/utils`: Contains utility functions and constants.


## Building the Project

To build the project, use the provided Makefile:

```sh
make build
```

This will create a binary named cr-product-1.0.0.bin.

## Contributing
Contributions are welcome. Please make sure to update tests as appropriate.

## License
MIT