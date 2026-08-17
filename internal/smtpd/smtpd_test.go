package smtpd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/emersion/go-smtp"
)

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("SMTP_ADDR", "")
	t.Setenv("DOMAIN_LIST", "")
	t.Setenv("TMAIL_REPORT_URL", "http://127.0.0.1:3000/api/report")
	t.Setenv("REPORT_HMAC_SECRET", "")
	t.Setenv("REPORT_MAX_BODY_SIZE", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Addr != ":25" {
		t.Fatalf("unexpected SMTP_ADDR default: %q", cfg.Addr)
	}
	if cfg.ReportMaxBodySize != 268435456 {
		t.Fatalf("unexpected REPORT_MAX_BODY_SIZE default: %d", cfg.ReportMaxBodySize)
	}
}

func TestLoadConfigRequiresReportURL(t *testing.T) {
	t.Setenv("TMAIL_REPORT_URL", "")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewServer(cfg); err == nil {
		t.Fatal("expected missing TMAIL_REPORT_URL to fail")
	}
}

func TestServerAcceptsOneRecipientPerTransaction(t *testing.T) {
	server, err := NewServer(Config{
		Addr:              ":25",
		ReportURL:         "http://127.0.0.1:3000/api/report",
		ReportMaxBodySize: 268435456,
	})
	if err != nil {
		t.Fatal(err)
	}
	if server.MaxRecipients != 1 {
		t.Fatalf("unexpected maximum recipients: %d", server.MaxRecipients)
	}
}

func TestRecipientDomainFilter(t *testing.T) {
	b := &backend{domains: map[string]struct{}{"example.com": {}}}
	s := &session{backend: b}
	if err := s.Rcpt("User@EXAMPLE.COM", nil); err != nil {
		t.Fatalf("expected configured domain to be accepted: %v", err)
	}
	if got := s.recipient; got != "User@example.com" {
		t.Fatalf("unexpected normalized recipient: %q", got)
	}
	err := s.Rcpt("user@invalid.example", nil)
	smtpErr, ok := err.(*smtp.SMTPError)
	if !ok || smtpErr.Code != 550 {
		t.Fatalf("expected SMTP 550, got %v", err)
	}
}

func TestRecipientWithoutDomainFilter(t *testing.T) {
	s := &session{backend: &backend{domains: map[string]struct{}{}}}
	if err := s.Rcpt("user@any.example", nil); err != nil {
		t.Fatalf("expected recipient to be accepted without DOMAIN_LIST: %v", err)
	}
}

func TestDataReportsRawMessageWithHMAC(t *testing.T) {
	const (
		secret    = "test-secret"
		recipient = "user@example.com"
		message   = "From: sender@example.org\r\nTo: user@example.com\r\nSubject: test\r\n\r\nhello"
	)
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/report" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("to"); got != recipient {
			t.Errorf("unexpected recipient: %q", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if string(body) != message {
			t.Errorf("unexpected message body: %q", body)
		}
		timestamp := r.Header.Get("X-Tmail-Timestamp")
		value := timestamp + "\nPOST\n/api/report\n" + recipient
		mac := hmac.New(sha256.New, []byte(secret))
		_, _ = mac.Write([]byte(value))
		expected := hex.EncodeToString(mac.Sum(nil))
		if got := r.Header.Get("X-Tmail-Signature"); got != expected {
			t.Errorf("unexpected signature: %q", got)
		}
		_, _ = io.WriteString(w, `{"code":0}`)
	}))
	defer receiver.Close()

	reportURL, err := url.Parse(receiver.URL + "/api/report")
	if err != nil {
		t.Fatal(err)
	}
	b := &backend{
		reportURL: reportURL,
		secret:    secret,
		domains:   map[string]struct{}{},
		client:    receiver.Client(),
	}
	s := &session{backend: b, recipient: recipient}
	if err := s.Data(strings.NewReader(message)); err != nil {
		t.Fatal(err)
	}
}

func TestDataRejectsOversizedMessage(t *testing.T) {
	reportURL, err := url.Parse("http://tmail.test/api/report")
	if err != nil {
		t.Fatal(err)
	}
	s := &session{
		backend: &backend{
			reportURL: reportURL,
			client: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				_, err := io.Copy(io.Discard, req.Body)
				return nil, err
			})},
		},
		recipient: "user@example.com",
	}
	err = s.Data(errorReader{err: smtp.ErrDataTooLarge})
	if err != smtp.ErrDataTooLarge {
		t.Fatalf("expected ErrDataTooLarge, got %v", err)
	}
}

type errorReader struct {
	err error
}

func (r errorReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
