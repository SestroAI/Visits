FROM golang:1.8.3

RUN mkdir -p /app/src/github.com/SestroAI/Visits
ADD . /app/src/github.com/SestroAI/Visits

WORKDIR /app/src/github.com/SestroAI/Visits
ENV GOPATH /app

RUN go build -o /app/bin/sestro-visits .

CMD ["/app/bin/sestro-visits"]
