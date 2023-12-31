package client

import (
	"net/url"
	"testing"
	"time"
)

func TestRequest_generateKey(t *testing.T) {
	type args struct {
		baseURL string
		path    string
		query   url.Values
		body    interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test cosmos key without params",
			args: args{
				baseURL: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/cosmos",
				path:    "validators/list.json",
			},
			want: "ukpgy7t9m_vLHvyQL82smBoTov4=",
		},
		{
			name: "test cosmos key with params",
			args: args{
				baseURL: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/cosmos",
				path:    "validators/list.json",
				query:   url.Values{"address": {"TQZskDJJRGAHifeKoQ7wLey42iGvwp3"}, "visible": {"false"}},
			},
			want: "jkkaXhzkelj5l3WE_B57Q1IY0Qo=",
		},
		{name: "test tron key without params ",
			args: args{
				baseURL: "https://api.trongrid.io",
				path:    "wallet/getaccount",
			},
			want: "PIoOx2azFYta4KMAtt0lttrqquM=",
		},
		{name: "test tron key with params 1",
			args: args{
				baseURL: "https://api.trongrid.io",
				path:    "wallet/getaccount",
				body: struct {
					Address string `json:"address"`
					Visible bool   `json:"visible"`
				}{Address: "TQZskDJJRGAHifeKoQ7wLC4QDyB2iGvwp2", Visible: true},
			},
			want: "h0noiR5a4M_RGQBH7805sgGl_HE=",
		},
		{name: "test tron key with params 2",
			args: args{
				baseURL: "https://api.trongrid.io",
				path:    "wallet/getaccount",
				body: struct {
					Address string `json:"address"`
					Visible bool   `json:"visible"`
				}{Address: "TQZskDJJRGAHifeKoQ7wLey42iGvwp3", Visible: false},
			},
			want: "Admv3wAXHkirPi4SaIXimDgLbow=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{BaseURL: tt.args.baseURL}
			if got := r.generateKey(tt.args.path, tt.args.query, tt.args.body); got != tt.want {
				t.Errorf("generateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

type (
	args struct {
		baseURL string
		path    string
		query   url.Values
		result  interface{}
	}
	response struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Website     string `json:"website"`
	}
	test struct {
		name string
		args args
	}
)

func testCollection() []test {
	return []test{
		{
			name: "test cosmos key without params",
			args: args{
				baseURL: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/cosmos/",
				path:    "validators/list.json",
				result:  new([]response),
			},
		},
		{
			name: "test cosmos key with params",
			args: args{
				baseURL: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/cosmos/",
				path:    "validators/list.json",
				query:   url.Values{"address": {"TQZskDJJRGAHifeKoQ7wLey42iGvwp3"}, "visible": {"false"}},
				result:  new([]response),
			},
		},
	}
}

func TestRequest_GetWithCache(t *testing.T) {
	for _, tt := range testCollection() {
		t.Run(tt.name, func(t *testing.T) {
			r := InitClient(tt.args.baseURL, nil)
			if err := r.GetWithCache(tt.args.result, tt.args.path, tt.args.query, time.Duration(1*time.Second)); err != nil {
				t.Errorf("GetWithCache was failed for %v, error %v", tt.name, err)
			}

			key := r.generateKey(tt.args.path, tt.args.query, nil)

			_, ok := memoryCache.cache.Get(key)

			if !ok {
				t.Errorf("GetWithCache could not find cache for %v", tt.name)
			}
		})
	}
}
