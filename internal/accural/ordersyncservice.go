package accural

import (
	"context"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/Alheor/gophermart/internal/repository"
	"strconv"
	"time"
)

type SyncService struct {
	SyncChan chan string
}

var ss *SyncService

const (
	StatusProcessed = `PROCESSED`
	StatusInvalid   = `INVALID`
)

func Init() {

	ss = &SyncService{SyncChan: make(chan string)}

	go func() {

		for {
			s := <-ss.SyncChan

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			cancel()

			go syncOrder(ctx, s)

			time.Sleep(1 * time.Second)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	loadFromBD(ctx)
}

func Sync(orderID string) {
	ss.SyncChan <- orderID
}

func syncOrder(ctx context.Context, orderID string) {

	data, err := connector.getOrderData(orderID)
	if err != nil {
		logger.GetLogger().Error(`Order sync error: ` + err.Error())

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
		ss.SyncChan <- orderID
	}
}

func loadFromBD(ctx context.Context) {
	logger.GetLogger().Info(`Loading order from BD ...`)

	list, err := repository.GetOrderRepository().GetOrderForProcessing(ctx)
	if err != nil {
		panic(err)
	}

	for _, order := range list {
		Sync(order.Number)
	}

	logger.GetLogger().Info(`Loaded ` + strconv.Itoa(len(list)) + ` orders`)
}
