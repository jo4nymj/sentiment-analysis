FROM golang:1.13

ENV PORT 5002

RUN mkdir /app
WORKDIR /app/
COPY . .

RUN go mod download && \
    go install .

CMD ["code.sentiments"]