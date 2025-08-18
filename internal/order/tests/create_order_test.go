package tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	sw "github.com/mushroomyuan/gorder/common/client/order"
	_ "github.com/mushroomyuan/gorder/common/config"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	ctx    = context.Background()
	server = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
	client *sw.ClientWithResponses
)

func TestMain(m *testing.M) {
	before()
	m.Run()

}
func before() {
	log.Printf("server=%s", server)
}

func TestCreateOrder_invalidParams(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items:      nil,
	})

	assert.Equal(t, http.StatusOK, response.StatusCode())
	//assert.Equal(t, 0, response.JSON200.Errno)

}

func TestCreateOrder_success(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items: []sw.ItemWithQuantity{
			{
				Id:       "test-item-1",
				Quantity: 10,
			},
		},
	})
	t.Logf("json200=%+v", response.JSON200)
	assert.Equal(t, http.StatusOK, response.StatusCode())
}

func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrdersJSONRequestBody) *sw.PostCustomerCustomerIdOrdersResponse {
	t.Helper()
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
