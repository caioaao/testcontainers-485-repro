package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/sync/errgroup"
)

var requests = []testcontainers.GenericContainerRequest{
	{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{pgListeningPort},
			WaitingFor: wait.
				ForSQL(nat.Port(pgListeningPort), "postgres", func(_ string, port nat.Port) string {
					// we ignore the host here due to issues with ipv4/ipv6 port mapping
					return fmt.Sprintf("host=127.0.0.1 user=%s password=%s dbname=%s port=%s sslmode=disable timezone=UTC",
						pgUser, pgPassword, pgDBName, port.Port())
				}).
				WithQuery("SELECT 1").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second),
			Env: map[string]string{
				"POSTGRES_DB":       pgDBName,
				"POSTGRES_USER":     pgUser,
				"POSTGRES_PASSWORD": pgPassword,
			},
		},
		Started: true,
	},
}

func main() {
	rand.Seed(time.Now().Unix())

	numContainers := flag.Int("n", 100, "number of containers to run")
	flag.Parse()

	g := new(errgroup.Group)

	var containers []testcontainers.Container
	fails := 0

	for i := 0; i < *numContainers; i++ {
		i := i
		g.Go(func() error {
			container, err := testcontainers.GenericContainer(context.Background(), requests[rand.Int()%len(requests)])
			if err != nil {
				fails++
				fmt.Printf("failed to start container %d: %v\n", i, err.Error())
				return err
			}

			containers = append(containers, container)

			return nil
		})
	}
	err := g.Wait()

	for _, container := range containers {
		container.Terminate(context.Background())
	}

	if err != nil {
		fmt.Printf("%d containers failed to start, out of %d\n", fails, *numContainers)
		os.Exit(-1)
	}
}

const (
	pgListeningPort = "5432/tcp"
	pgUser          = "user"
	pgPassword      = "changeit"
	pgDBName        = "db"
)

const (
	gcpPubsubListeningPort = "8538/tcp"
)
