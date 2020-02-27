# httpmole

[![Build Status](https://travis-ci.com/jcchavezs/httpmole.svg?branch=master)](https://travis-ci.com/jcchavezs/httpmole)

**httpmole** provides a HTTP mock server that will act as a mole among your services, telling you everything http clients send to it and responding them whatever you want it to respond. Just like an actual mole.

It supports:

- `response-status` and `response-header` to quickly spin up a http server.
- `response-file` option where **you can modify the response in real time** using a file.
- `response-from` so it can act as a proxy and you can inspect the request/response going to a given service.

## Install

```bash
go install github.com/jcchavezs/httpmole/cmd/httpmole
```

## Usage

### Using the binary

```bash
httpmole -p=8082 -response-status=200
```

or using a response file:

```bash
httpmole -p=8082 -response-file=./myresponse.json
vim ./myresponse.json
```

```json
// myresponse.json
{
    "status_code": 200,
    "headers": {
        "content-type": "application/json"
    },
    "body": {
        "message": "I am real"
    }
}
```

or proxying a service to inspect the requests:

```bash
httpmole -p=8082 -response-from=therealservice:8082
```

### Using docker

```bash
docker run -p "8081:8081" -v `pwd`/response.json:/httpmole/response.json -response-file=/httpmole/response.json jcchavezs/httpmole
```

Docker image is [hosted in dockerhub](https://hub.docker.com/repository/docker/jcchavezs/httpmole
)
