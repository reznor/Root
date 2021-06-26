## Invocation

`$ go run main.go input.txt`

`$ cat input.txt | go run main.go`

## Notes

### Directory layout

Having the `eventhandlers` intermediate directory allows for adding common shared helpers/utilities there in the future; each eventhandler getting its own subdirectory in there leads to stronger ACL isolation across different `eventhandlers`, with the option of completely different test patterns and build strategies, if so desired.