package nats

import (
    "github.com/nats-io/nats.go"
)

func Connect(url string) (*nats.Conn, error) {
    return nats.Connect(url)
}
