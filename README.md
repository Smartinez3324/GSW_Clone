# Launch Ground Software
Ground software code for RIT Launch Initiative. Responsible for receiving telemetry over Ethernet and publishing to shared memory for other applications to use.

### Compiling
By Golang convention, all files meant to be executed are stored in the cmd folder. 
* The GSW service can be compiled by running `go build cmd/gsw_service.go` from the project root directory
* GSW applications are stored in subdirectories within the cmd folder and the `go build` command can be ran on go files within those subdirectories

### Running
You can always run the GSW service by doing a `./gsw_service` after building. For running any Go program though, instead of doing `go build (FILE_PATH)` you can do `go run (FILE_PATH` instead.
(TODO) Running as a service

### Unit Tests
There are several unit tests that can be ran. You can do a `go test ./...` from the root project directory to execute all tests. It is also recommended to run with the -cover
flag to get coverage statements.
