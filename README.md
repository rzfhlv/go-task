# Just go-task

## Requirement

- go version go1.23.4 or higher
- docker and docker-compose

## How to use

- clone from repo: 

    ``` git clone https://github.com/rzfhlv/go-task.git ```

- in your root project directory copy environment file:

    ``` cp config.yml.tmpl config.yml ```

- golang initialize:

    ``` go mod tidy ```

- run the infrastructure:

    ``` make deps-up ```

- run the migration:

    ``` make migrate-up ```

- run the app:

    ``` make run ```

- application running on port 8080 by default

- postaman colletion available on docs directory