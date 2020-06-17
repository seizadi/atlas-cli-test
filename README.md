# test

This generated README.md file loosely follows a [popular template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)._

One paragraph of project description goes here.

## Getting Started

If you are here you created the app using atlas-cli
```bash
atlas init-app -name=test -expand=expand.txt -db=true -registry=<your registry> -gateway -health -helm
```

There is an [issue with atlas-cli](https://github.com/infobloxopen/atlas-cli/issues/60).
The gentool which is used to generate the code is trailing the version used
by atlas-cli, result is that there
is an older version of grpc-gateway code gen versus go module declarations.
When you try to compile you get errors, for now we patch go.mod:
```bash
diff --git a/go.mod b/go.mod
index 46adf47..0c3ad71 100644
--- a/go.mod
+++ b/go.mod
@@ -6,7 +6,7 @@ require (
        github.com/golang/protobuf v1.4.2
        github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
        github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
-       github.com/grpc-ecosystem/grpc-gateway v1.14.6
+       github.com/grpc-ecosystem/grpc-gateway v1.9.6
        github.com/infobloxopen/atlas-app-toolkit v0.21.1
```

It should compile:

```bash
go mod vendor
go build -o test cmd/server/*.go
```

You should be able to run the full make:
```bash
make all
```
You will see a successful make and in my case an image built for my registry:
```bash
Successfully tagged soheileizadi/test:e07806a-unsupported
```

Lets start running it local, this assumes you have postgres database running on
localhost listening to 5432, we do this by running a docker postgres image:
```bash
docker run -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=atlas_cli_test -p 5432:5432 postgres
```

Make sure it is running
```bash
❯ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
1eee74c1cf18        postgres            "docker-entrypoint.s…"   27 seconds ago      Up 26 seconds       0.0.0.0:5432->5432/tcp   objective_leavitt
```

Put a rule in make file to run migration locally using migrate tool:
```bash
DATABASE_HOST=localhost:5432
DATABASE_NAME=test
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

migrate:
    @migrate -database 'postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST)/$(DATABASE_NAME)?sslmode=disable' -path ./db/migration up
```

Atlas CLI is missing a couple of migration files you can for now copy them from here:
[00001_set_updated_at_func.down.sql](https://github.com/seizadi/cmdb/blob/master/db/migrations/00001_set_updated_at_func.down.sql)
[00001_set_updated_at_func.up.sql](https://github.com/seizadi/cmdb/blob/master/db/migrations/00001_set_updated_at_func.up.sql)

Migration should work now:
```bash
❯ make migrate
migrate -database 'postgres://localhost:5432/test?sslmode=disable' -path ./db/migration up
1/u set_updated_at_func (7.167753ms)
2/u accounts (14.064373ms)
3/u users (20.039864ms)
4/u groups (25.608444ms)
```

You can now run the app:
```bash
❯ ./test
...
 ./test
WARNING: Package "github.com/golang/protobuf/protoc-gen-go/generator" is deprecated.
        A future release of golang/protobuf will delete this package,
        which has long been excluded from the compatibility promise.

2020/06/16 23:04:17 Serving from default values, environment variables, and/or flags
{"level":"debug","msg":"serving internal http at \"0.0.0.0:8081\"","time":"2020-06-16T23:04:17-07:00"}
{"level":"info","msg":"serving gRPC at 0.0.0.0:9090","time":"2020-06-16T23:04:17-07:00"}
{"level":"info","msg":"serving http at 0.0.0.0:8080","time":"2020-06-16T23:04:17-07:00"}
```
Now we get version output but other endpoints are not working!
```bash
❯ curl http://localhost:8080/atlas-cli-test/v1/version | jq
{
  "version": "0.0.1"
}

❯ curl -H "Authorization: Bearer $JWT" http://localhost:8080/atlas-cli-test/v1/accounts | jq
{
  "error": [
    {
      "message": "Not Implemented"
    }
  ]
}
```

This is a problem with endpoint registeration, you need to modify main.go and replace the 
[following commit](https://github.com/seizadi/atlas-cli-test/commit/ef6b060a42756b990a2f3840f59a6bf9bbf7dd9c)
and changes to:
   * grpc.go
   * main.go
   * servers.go
   * endpoints.go

After these changes you should be reach the accounts endpoint and others:
```bash
curl http://localhost:8080/atlas-cli-test/v1/accounts | jq
{}
```
Now you can add an Account:
```bash
curl http://localhost:8080/atlas-cli-test/v1/accounts -d '{"name": "acme", "description": "Acme Corp"}' | jq
{
  "result": {
    "id": "atlas-cli-test/accounts/1",
    "name": "acme",
    "description": "Acme Corp"
  }
}

curl -X PATCH http://localhost:8080/atlas-cli-test/v1/accounts/1 -d '{"description": "Acme Inc"}' | jq
{
  "result": {
    "id": "atlas-cli-test/accounts/1",
    "name": "acme",
    "description": "Acme Inc"
  }
}

curl -X PUT http://localhost:8080/atlas-cli-test/v1/accounts/1 -d '{"name": "Acme New"}' | jq
{
  "result": {
    "id": "atlas-cli-test/accounts/1",
    "name": "Acme New"
  }
}

curl http://localhost:8080/atlas-cli-test/v1/accounts | jq
{
  "results": [
    {
      "id": "atlas-cli-test/accounts/1",
      "name": "Acme New"
    }
  ]
}

curl -X DELETE http://localhost:8080/atlas-cli-test/v1/accounts/1 | jq
{}
```
To view swagger:
```bash
 make run-swagger-ui 
```

- EXTRA CREDIT setup some relationships

```bash
curl http://localhost:8080/atlas-cli-test/v1/version | jq
export JWT="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50SUQiOjF9.GsXyFDDARjXe1t9DPo2LIBKHEal3O7t3vLI3edA7dGU"
curl -H "Authorization: Bearer $JWT" http://localhost:8080/atlas-cli-test/v1/accounts
curl -H "Authorization: Bearer $JWT" http://localhost:8080/atlas-cli-test/v1/users
curl -H "Authorization: Bearer $JWT" http://localhost:8080/atlas-cli-test/v1/groups
```

### Prerequisites

What things you need to install the software and how to install them.

```
Give examples
```

### Installing

A step-by-step series of examples that tell you have to get a development environment running.

Say what the step will be.

```
Give the example
```

And repeat.

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo.



## Deployment

Add additional notes about how to deploy this application. Maybe list some common pitfalls or debugging strategies.

## Running the tests

Explain how to run the automated tests for this system.

```
Give an example
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/apptriton/test/tags).
