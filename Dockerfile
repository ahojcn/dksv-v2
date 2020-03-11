FROM library/golang

# Godep for vendoring
COPY ["./", "/go/src/"]

RUN cd /go/src/github.com/beego/bee && go install

ENV PATH $PATH:$GOPATH/bin

# RUN cd /go/src/dksv-v2 && bee run

CMD cd /go/src/dksv-v2 && bee run

EXPOSE 8080
