package coinank

import (
	"context"
	"encoding/json"
	"nofx/provider/coinank/coinank_enum"
	"testing"
)

func TestGetLastPrice(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.GetLastPrice(context.TODO(), "BTCUSDT", "Binance", "SWAP")
	if err != nil {
		t.Fatal(err)
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}

func TestGetCoinMarketCap(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.GetCoinMarketCap(context.TODO(), "BTC")
	if err != nil {
		t.Fatal(err)
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}
