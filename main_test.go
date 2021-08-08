package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	containesvr "testcontainers-go-e2e/containersvr"
	"testcontainers-go-e2e/db"
	"testing"

	_ "github.com/lib/pq"
	tc "github.com/testcontainers/testcontainers-go"
	"gotest.tools/assert"
)

var ctx = context.Background()
var _userServiceURL, _ticketServiceURL string

func TestMain(m *testing.M) {
	// os.Chdir("..") // back to the parent directory of current directory

	var network = tc.NetworkRequest{
		Name:   "integration-test-network",
		Driver: "bridge",
	}

	provider, err := tc.NewDockerProvider()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := provider.GetNetwork(ctx, network); err != nil {
		if _, err := provider.CreateNetwork(ctx, network); err != nil {
			log.Fatal(err)
		}
	}

	postgresConfig := db.PostgresConfig{
		Password: "password",
		User:     "postgres",
		DB:       "userservice",
		Port:     "5432/tcp",
	}

	postgresInternal, mappedPostgres := postgresConfig.
		StartContainer(network.Name)
	log.Println("postgres running at: ", mappedPostgres)

	internalUser, mappedUser := containesvr.UserServiceConfig{PostgresURL: postgresInternal, Port: "8080/tcp"}.
		StartContainer(network.Name)
	log.Println("user service running at: ", mappedUser)

	_, _ticketServiceURL = containesvr.TicketServiceConfig{UserServiceURL: internalUser, Port: "8080/tcp"}.
		StartContainer(network.Name)
	log.Println("ticket service running at: ", _ticketServiceURL)

	_userServiceURL = mappedUser

	os.Exit(m.Run())
}

func Test_Integrations(t *testing.T) {
	var createdUser User
	var theUser User
	t.Run("create", func(t *testing.T) {
		resp, _ := http.Post(
			_userServiceURL+"/users",
			"application/json",
			JsonReader(User{Name: "berliner", ID: '5'}))

		ReadResp(resp, &createdUser)
		assert.Equal(t, "berliner", createdUser.Name)

		resp, error := http.Get(fmt.Sprintf("%s/users/%d", _userServiceURL, createdUser.ID))
		if error != nil {
			log.Println(error)
			t.Fatal()
		}
		ReadResp(resp, &theUser)
		assert.Equal(t, theUser, createdUser)
	})
	var ticket Ticket
	var userTicket Ticket
	t.Run("create a ticket", func(t *testing.T) {
		resp, _ := http.Post(
			_ticketServiceURL+"/tickets",
			"application/json",
			JsonReader(TicketPost{Movie: "Berlin Syndrome", UserID: createdUser.ID}))
		ReadResp(resp, &ticket)

		resp, _ = http.Get(fmt.Sprintf("%s/tickets/%s", _ticketServiceURL, ticket.ID))
		ReadResp(resp, &userTicket)

		assert.Equal(t, createdUser, userTicket.User)
		assert.Equal(t, ticket.ID, userTicket.ID)
		assert.Equal(t, "Berlin Syndrome", userTicket.Movie)
	})
}
