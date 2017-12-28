FROM golang:1.9

ADD ./ /go/src/github.com/whywaita/slack-to-fluentd/logger/
WORKDIR /go/src/github.com/whywaita/slack-to-fluentd/logger/

RUN go get -u github.com/golang/dep/cmd/dep \
  && dep ensure \
  && go build -o logger main.go

CMD ["/go/src/github.com/whywaita/slack-to-fluentd/logger/logger"]
