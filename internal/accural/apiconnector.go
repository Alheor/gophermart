package accural

import (
	"encoding/json"
	"errors"
	"github.com/Alheor/gophermart/internal/config"
	"github.com/Alheor/gophermart/internal/entity"
	"io"
	"net/http"
	"time"
)

const (
	apiGetOrderDataPath = `/api/orders/`
)

type API interface {
	getOrderData(orderID string) error
}

type APIConnector struct {
	client *http.Client
}

var connector *APIConnector

func init() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}

	connector = new(APIConnector)
	connector.client = &http.Client{Transport: tr}
}

func (ac *APIConnector) getOrderData(orderID string) (*entity.AccrualOrder, error) {

	resp, err := connector.client.Get(config.Options.AccrualSystemAddress + apiGetOrderDataPath + orderID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(``)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var el entity.AccrualOrder
	err = json.Unmarshal(body, &el)
	if err != nil {
		return nil, err
	}

	return &el, nil
}
