FROM denismakogon/gocv-alpine:4.0.1-buildstage as build-stage

RUN go get -u -d gocv.io/x/gocv
RUN cd $GOPATH/src/gocv.io/x/gocv && go build -o $GOPATH/bin/gocv-version ./cmd/version/main.go
WORKDIR $GOPATH/src/gocv-playground
