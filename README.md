# httpmole

**httpmole** provides a mock server that will lie to http clients, saying as response whatever you tell it to say.

It provides support for a `response-file` that you can edit on real-time.

## Usage

### Using golang

```bash
make build // builds the binary
./httpmole -p=8082 -response-status=200
```

or reading a response file:

```bash
make build
./httpmole -p=8082 -response-file=./myresponse.json
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
