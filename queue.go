package main

import (
	"log"

	"github.com/sivamgr/kstreamdb"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

var tickChannel chan kiteticker.Tick
var ksocket kstreamdb.Socket

func setupTickHandler() {
	const channelCapacity = 256
	tickChannel = make(chan kiteticker.Tick, channelCapacity)
	var e error
	ksocket, e = kstreamdb.StartStreaming(AppConfig.DataManagement.PublishSocket)
	if e == nil {
		db.RecordStream(&ksocket)
	}
	go handleKiteTicks()
}

func handleKiteTicks() {
	for {
		if msg, ok := <-tickChannel; ok {
			tick := buildTick(&msg)
			ksocket.Publish(tick)
		} else {
			log.Println("Channel error with handleKiteTicks")
			break
		}
	}

}

func buildDepth(d [5]kiteticker.DepthItem) [5]kstreamdb.DepthItem {
	var kd [5]kstreamdb.DepthItem
	for i := 0; i < 5; i++ {
		kd[i].Price = float32(d[i].Price)
		kd[i].Quantity = d[i].Quantity
		kd[i].Orders = d[i].Orders
	}
	return kd
}

func buildTick(k *kiteticker.Tick) kstreamdb.TickData {
	t := kstreamdb.TickData{
		TradingSymbol: mapTokenToSymbol[k.InstrumentToken],
		IsTradable:    k.IsTradable,

		Timestamp: k.Timestamp.Time,

		LastTradeTime:      k.LastTradeTime.Time,
		LastPrice:          float32(k.LastPrice),
		LastTradedQuantity: k.LastTradedQuantity,

		AverageTradePrice: float32(k.AverageTradePrice),

		VolumeTraded:      k.VolumeTraded,
		TotalBuyQuantity:  k.TotalBuyQuantity,
		TotalSellQuantity: k.TotalSellQuantity,

		DayOpen:      float32(k.OHLC.Open),
		DayHighPrice: float32(k.OHLC.High),
		DayLowPrice:  float32(k.OHLC.Low),
		LastDayClose: float32(k.OHLC.Close),

		OI:        k.OI,
		OIDayHigh: k.OIDayHigh,
		OIDayLow:  k.OIDayLow,

		Bid: buildDepth(k.Depth.Buy),
		Ask: buildDepth(k.Depth.Sell),
	}
	return t
}
