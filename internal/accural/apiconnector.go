package accural

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Alheor/gophermart/internal/models"
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
var serviceAddr string

func InitConnector(accrualAddr string) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}

	connector = new(APIConnector)
	connector.client = &http.Client{Transport: tr}

	serviceAddr = accrualAddr
}

func (ac *APIConnector) getOrderData(orderID string) (*models.AccrualOrder, error) {

	resp, err := connector.client.Get(serviceAddr + apiGetOrderDataPath + orderID)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(`accural response status code: ` + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var el models.AccrualOrder
	err = json.Unmarshal(body, &el)
	if err != nil {
		return nil, err
	}

	return &el, nil
}
