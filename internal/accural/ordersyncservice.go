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

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			println(s)
			go syncOrder(ctx, s)
			time.Sleep(1 * time.Second)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	loadFromBD(ctx)
}

func Sync(orderId string) {
	ss.SyncChan <- orderId
}

func syncOrder(ctx context.Context, orderId string) {
	data, err := connector.getOrderData(orderId)

	if err != nil {
		logger.GetLogger().Error(`Order sync error: ` + err.Error())
		ss.SyncChan <- orderId
		return
	}

	if data.Status != StatusProcessed && data.Status != StatusInvalid {
		ss.SyncChan <- orderId
	}

	err = repository.GetOrderRepository().ChangeOrder(ctx, data)
	if err != nil {
		panic(err)
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
