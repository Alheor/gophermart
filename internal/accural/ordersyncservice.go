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

func InitService(ctx context.Context) {

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
			logger.Error(`Chan error: `, err)
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
		logger.Error(`Order sync error: `, err)

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
		logger.Error(`Change order error: `, err)

		ss.ErrChan <- err
		ss.SyncChan <- orderID
	}
}

func loadFromBD(ctx context.Context) {
	logger.Info(`Loading order from BD ...`)

	list, err := repository.GetOrderRepository().GetOrderForProcessing(ctx)
	if err != nil {
		logger.Error(`Error load orders from DB: `, err)

		ss.ErrChan <- err
		return
	}

	for _, order := range list {
		Sync(order.Number)
	}

	logger.Info(`Loaded ` + strconv.Itoa(len(list)) + ` orders`)
}
