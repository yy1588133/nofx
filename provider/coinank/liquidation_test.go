package coinank

import (
	"context"
	"encoding/json"
	"nofx/provider/coinank/coinank_enum"
	"testing"
	"time"
)

func TestLiquidationExchangeStatistics(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.LiquidationExchangeStatistics(context.TODO(), "BTC")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total <= 0 {
		t.Error("total amount is negative")
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}

func TestLiquidationCoinAggHistory(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.LiquidationCoinAggHistory(
		context.TODO(),
		"BTC",
		coinank_enum.Hour1,
		time.Now().UnixMilli(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) == 0 {
		t.Fatal("empty response")
	}
	if resp[0].All.LongTurnover <= 0 {
		t.Error("longTurnover is negative")
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}

func TestLiquidationHistory(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.LiquidationHistory(
		context.TODO(),
		coinank_enum.Binance,
		"BTCUSDT",
		coinank_enum.Hour1,
		time.Now().UnixMilli(),
		10,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) == 0 {
		t.Fatal("empty response")
	}
	if resp[0].LongTurnover <= 0 {
		t.Error("longTurnover is negative")
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}

func TestLiquidationOrders(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.LiquidationOrders(
		context.TODO(),
		"BTC",
		coinank_enum.Binance,
		"long",
		1000,
		time.Now().UnixMilli(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) == 0 {
		t.Fatal("empty response")
	}
	if resp[0].Price <= 0 {
		t.Error("price is negative")
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}

func TestLiquidationOrdersNoArgs(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.LiquidationOrders(context.TODO(), "", "", "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) == 0 {
		t.Fatal("empty response")
	}
	if resp[0].Price <= 0 {
		t.Error("price is negative")
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}
