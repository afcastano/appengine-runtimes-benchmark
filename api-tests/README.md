See https://blog.narenarya.in/concurrent-http-in-go.html for examples.

In the root folder:
```
export GOPATH=$PWD
export GOBIN=$GOPATH/bin
PATH=$PATH:$GOPATH:$GOBIN
export PATH=$PATH
```

Run
```
go build concurrent.go 
```
