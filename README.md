# httpmole

**httpmole** provides a HTTP mock server that will act as a mole among your services, telling you everything http clients send to it and telling them whatever you want. Just like an actual mole.

It provides support for a `response-file` option where **you can modify the response in real time**.

## Install

```bash
go install github.com/jcchavezs/httpmole
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

### Using docker

```bash
docker run -p "8081:8081" -v `pwd`/response.json:/httpmole/response.json -response-file=/httpmole/response.json
```
