# Hexagonal Architecture 
## About
This repository is for a demo service I used to make a presentation on hexagonal architecture at HelloFresh on 31.01.2023.
The service is a url shortener.
## Goal
To demonstrate how to use hexagonal architecture to build a go service.
## Features
Shortens a long http url
## API Documentation
None
## Tools
Tools used during the development of the service are;
- [Go](https://go.dev/) - Golang
- [Postgresql](https://www.postgresql.org/) - this is a database server
## Requirements
- Go 1.17+.
- Postgres 12+
## Tests
Even God commands us to run tests: 1 Thessalonians 5:21; "Test all things."
So to run tests, go to your command line prompt and execute the following command
```sh
   $ make test-unit
```
## Running the application
To run this application, execute the following command.
```sh
    $ make dev-up
    $ make dev-ssh
    $ make migrations
    $ shortener shorten https://www.test.com
```
## Presentation slides
You can access the power point slide [here](https://docs.google.com/presentation/d/1isvtaIdXbn67voh6uiWi_1tee8nX90Mfsxv0BCzeLBg/edit?usp=sharing)
## Homework
As homework fork this repository and add a feature to retrieve the long http url using the short url. 