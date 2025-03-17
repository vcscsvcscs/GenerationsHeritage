package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck_GetStatus(t *testing.T) {
	var hc HealthCheck

	type fields struct {
		status string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test GetStatus with status ok",
			fields: fields{
				status: "ok",
			},
			want: "ok",
		},
		{
			name: "Test GetStatus with status nok",
			fields: fields{
				status: "nok",
			},
			want: "nok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc = &healthCheck{
				status: tt.fields.status,
			}
			if got := hc.GetStatus(); got != tt.want {
				t.Errorf("HealthCheck.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthCheck_SetStatus(t *testing.T) {
	hc := New()

	type args struct {
		status string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test SetStatus with status ok",
			args: args{
				status: "ok",
			},
			want: "ok",
		},
		{
			name: "Test SetStatus with status nok",
			args: args{
				status: "nok",
			},
			want: "nok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc.SetStatus(tt.args.status)
			if got := hc.GetStatus(); got != tt.want {
				t.Errorf("HealthCheck.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthCheck_HealthCheckHandler(t *testing.T) {
	r := gin.Default()
	hc := New()
	r.GET("/health", hc.HealthCheckHandler())
	type args struct {
		status string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		statusCode int
	}{
		{
			name: "Test respond with status ok",
			args: args{
				status: "ok",
			},
			want:       `{"status":"ok"}`,
			statusCode: http.StatusOK,
		},
		{
			name: "Test respond with status nok",
			args: args{
				status: "nok",
			},
			want:       `{"status":"nok"}`,
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc.SetStatus(tt.args.status)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/health", nil)
			r.ServeHTTP(w, req)

			if got := w.Code; got != tt.statusCode {
				t.Errorf("HealthCheck response status code = %v, want %v", got, tt.statusCode)
			}
			if got := w.Body.String(); got != tt.want {
				t.Errorf("HealthCheck response = %v, want %v", got, tt.want)
			}
		})
	}
}
