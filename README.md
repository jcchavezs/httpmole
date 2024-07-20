# httpmole

[![CI](https://github.com/jcchavezs/httpmole/actions/workflows/ci.yaml/badge.svg)](https://github.com/jcchavezs/httpmole/actions/workflows/ci.yaml)

**httpmole** provides a HTTP mock server that will act as a mole among your services, telling you everything http clients send to it and responding them whatever you want it to respond. Just like an actual mole.

<p align="center">
  <img width="640" height="356" src="images/screencast.gif">
</p>

**Features:**

- Use `response-status` and `response-header` to quickly spin up a http server.
- Use `response-file` to **modify the response in real time** using a text editor.
- Use `response-from` to act as a proxy and be able to inspect the request/response going to a given service.
- Cascade several calls using `/proxy` reserved endpoint.

## Install

```bash
go install github.com/jcchavezs/httpmole/cmd/httpmole@latest
```

## Usage

### Using the binary

```bash
httpmole -p=8082 -response-status=200
```

or using a response file:

```bash
httpmole -p=8082 -response-file=./myresponse.json
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

or proxying a service to inspect the incoming requests:

```bash
httpmole -p=8082 -response-from=therealservice:8082
```

or proxying a cascaring call to different services to simulate a distributed transaction:

```bash
# terminal 1
httpmole -p=8081 -response-status=201

# terminal 2
httpmole -p=8082 -response-status=202

# terminal 3
httpmole -p=8083 -response-status=203
```

and running

```bash
$ curl -i http://localhost:8081/proxy/localhost:8082/proxy/localhost:8083

HTTP/1.1 203 Non-Authoritative Information
Content-Length: 0
Date: Wed, 17 Jul 2024 08:51:00 GMT
Server-Timing: app;dur=0.00
```

### Using docker

```bash
docker run -p "10080:10080" ghcr.io/jcchavezs/httpmole -response-status=201
```

or pass a response file over volumes

```bash
docker run -p "10080:10080" -v `pwd`/response.json:/httpmole/response.json ghcr.io/jcchavezs/httpmole -response-file=/httpmole/response.json
```

**httpmole** is heavily inspired by [httplab](https://github.com/gchaincl/httplab)
