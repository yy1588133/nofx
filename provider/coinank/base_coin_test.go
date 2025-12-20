package coinank

import (
	"context"
	"encoding/json"
	"nofx/provider/coinank/coinank_enum"
	"testing"
)

func TestListCoin(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.ListCoin(context.TODO(), "SPOT")
	if err != nil {
		t.Fatal(err)
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}

func TestListSymbols(t *testing.T) {
	client := NewCoinankClient(coinank_enum.MainUrl, requireTestApikey(t))
	resp, err := client.ListSymbols(context.TODO(), "Binance", "SWAP")
	if err != nil {
		t.Fatal(err)
	}
	res, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}
