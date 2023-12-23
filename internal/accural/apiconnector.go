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

type Api interface {
	getOrderData(orderId string) error
}

type ApiConnector struct {
	client *http.Client
}

var connector *ApiConnector

func init() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}

	connector = new(ApiConnector)
	connector.client = &http.Client{Transport: tr}
}

func (ac *ApiConnector) getOrderData(orderId string) (*entity.AccrualOrder, error) {

	//body := []byte(`{
	//   "order": "` + orderId + `",
	//   "status": "PROCESSED",
	//	"accrual": 1
	//}`)

	resp, err := connector.client.Get(config.Options.AccrualSystemAddress + apiGetOrderDataPath + orderId)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(``)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var el entity.AccrualOrder
	err = json.Unmarshal(body, &el)
	if err != nil {
		return nil, err
	}

	return &el, nil
}
