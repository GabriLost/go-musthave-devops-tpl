package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	defaultTestPort = "8081"
	tcp             = "tcp"
)

func Test_SendCounter_And_SendGauge(t *testing.T) {
	type fields struct {
		name  string
		value int64
	}
	tests := []struct {
		name        string
		fields      fields
		client      *http.Client
		want        bool
		expectError bool
	}{
		{
			name:        "Send Counter",
			fields:      fields{name: "PollCount", value: 9},
			client:      &http.Client{Timeout: defaultTimeout},
			want:        true,
			expectError: false,
		},
		{
			name:        "Send Gauge",
			fields:      fields{name: "Alloc", value: 1000000},
			client:      &http.Client{Timeout: defaultTimeout},
			want:        true,
			expectError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := counter{
				name:  test.fields.name,
				value: test.fields.value,
			}

			l, err := net.Listen(tcp, defaultServer+":"+defaultTestPort)
			if err != nil {
				log.Fatal(err)
			}
			serveMux := http.NewServeMux()
			serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			srv := httptest.NewUnstartedServer(serveMux)
			err = srv.Listener.Close()
			if err != nil {
				return
			}
			srv.Listener = l
			srv.Start()

			defer srv.Close()

			got, err := c.SendCounter(test.client)
			if (err != nil) != test.expectError {
				t.Errorf("counter.SendCounter() error = %v, expectError %v", err, test.expectError)
				return
			}
			if got != test.want {
				t.Errorf("counter.SendCounter() = %v, want %v", got, test.want)
			}
		})
	}
}
