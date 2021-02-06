# Web micro-service

My assumptions are:

1. This service is placed after another one that manages authentication and authorization.
2. Because of that, we use HTTP instead of HTTPS.
3. This service is called from other services, not by users directly, therefore we trust the input information is valid.

Finally, I have chosen Gin as server module because when comparing with others looks more developed and with more
documentation. Moreover, I have also chosen the ArangoDB because I want to use it for a personal project and it has a go
driver.

## What this service does

This service does the following tasks:

1. Validates the input has the correct format, which is:
    ```json
    {
      "userId": "string",
      "username": "string"
    }
    ```
2. Adds a register to the DB with the `userId` field, and a timestamp to track when this service is called.
3. Returns a response with a hello message:
    ```json
    {
      "message": "Hello <username>!"
    }
    ```

## Deployment

First you need to deploy the database using the following command in `./docker` folder:

```shell
docker-compose up -d
```

Then you have to start the app at `./`:

```shell
go run .
```

Finally, test the server with:

```shell
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"userId": "123", "username":"Julio"}' \
    http://localhost:8080/echo
```

> *Note*: you can observe the DB through its interactive page at http://localhost:8529. Login using "root" for the username and "test" for the password.
