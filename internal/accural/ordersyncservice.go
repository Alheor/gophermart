package accural

import (
	"context"
	"strconv"
	"time"

	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/repository"
)

type SyncService struct {
	SyncChan chan string
	ErrChan  chan error
}

var ss *SyncService

const (
	StatusProcessed = `PROCESSED`
	StatusInvalid   = `INVALID`
)

func Init(ctx context.Context) {

	InitConnector()

	ss = &SyncService{SyncChan: make(chan string), ErrChan: make(chan error)}

	go func() {
		for {
			s := <-ss.SyncChan

			go syncOrder(ctx, s)

			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			err := <-ss.ErrChan
			handleError(err)
		}
	}()

	loadFromBD(ctx)
}

func Sync(orderID string) {
	ss.SyncChan <- orderID
}

func syncOrder(ctx context.Context, orderID string) {

	data, err := connector.getOrderData(orderID)
	if err != nil {
		logger.GetLogger().Error(`Order sync error: ` + err.Error())
		ss.ErrChan <- err
		ss.SyncChan <- orderID
		return
	}

	if data == nil {
		return
	}

	if data.Status != StatusProcessed && data.Status != StatusInvalid {
		ss.SyncChan <- orderID
	}

	err = repository.GetOrderRepository().ChangeOrder(ctx, data)
	if err != nil {
		logger.GetLogger().Error(`Change order error: ` + err.Error())
		ss.ErrChan <- err
		ss.SyncChan <- orderID
	}
}

func handleError(err error) {
	logger.GetLogger().Error(err.Error())
}

func loadFromBD(ctx context.Context) {
	logger.GetLogger().Info(`Loading order from BD ...`)

	list, err := repository.GetOrderRepository().GetOrderForProcessing(ctx)
	if err != nil {
		logger.GetLogger().Error(`Error load orders from DB: ` + err.Error())
		ss.ErrChan <- err
		return
	}

	for _, order := range list {
		Sync(order.Number)
	}

	logger.GetLogger().Info(`Loaded ` + strconv.Itoa(len(list)) + ` orders`)
}
