package smtpd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/emersion/go-smtp"
)

const (
	maxRecipients    = 1
	operationTimeout = 5 * time.Minute
)

type Config struct {
	Addr              string   `env:"SMTP_ADDR" envDefault:":25"`
	DomainList        []string `env:"DOMAIN_LIST"`
	ReportURL         string   `env:"TMAIL_REPORT_URL"`
	ReportSecret      string   `env:"REPORT_HMAC_SECRET"`
	ReportMaxBodySize int64    `env:"REPORT_MAX_BODY_SIZE" envDefault:"268435456"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg Config) validate() error {
	if cfg.Addr == "" {
		return errors.New("SMTP_ADDR must not be empty")
	}
	if cfg.ReportURL == "" {
		return errors.New("TMAIL_REPORT_URL is required")
	}
	if cfg.ReportMaxBodySize <= 0 {
		return fmt.Errorf("invalid REPORT_MAX_BODY_SIZE: %d", cfg.ReportMaxBodySize)
	}
	return nil
}

func NewServer(cfg Config) (*smtp.Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	reportURL, err := url.ParseRequestURI(cfg.ReportURL)
	if err != nil || (reportURL.Scheme != "http" && reportURL.Scheme != "https") || reportURL.Host == "" {
		return nil, fmt.Errorf("invalid TMAIL_REPORT_URL: %q", cfg.ReportURL)
	}
	if reportURL.RawQuery != "" || reportURL.Fragment != "" {
		return nil, errors.New("TMAIL_REPORT_URL must not contain a query or fragment")
	}

	domains := make(map[string]struct{}, len(cfg.DomainList))
	for _, domain := range cfg.DomainList {
		domain = normalizeDomain(domain)
		if domain != "" {
			domains[domain] = struct{}{}
		}
	}
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "localhost"
	}
	backend := &backend{
		reportURL: reportURL,
		secret:    cfg.ReportSecret,
		domains:   domains,
		client:    &http.Client{Timeout: operationTimeout},
	}
	server := smtp.NewServer(backend)
	server.Addr = cfg.Addr
	server.Domain = hostname
	server.MaxMessageBytes = cfg.ReportMaxBodySize
	server.MaxRecipients = maxRecipients
	server.ReadTimeout = operationTimeout
	server.WriteTimeout = operationTimeout
	return server, nil
}

type backend struct {
	reportURL *url.URL
	secret    string
	domains   map[string]struct{}
	client    *http.Client
}

func (b *backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &session{backend: b}, nil
}

type session struct {
	backend   *backend
	recipient string
}

func (s *session) Mail(_ string, _ *smtp.MailOptions) error {
	return nil
}

func (s *session) Rcpt(to string, _ *smtp.RcptOptions) error {
	address, domain, err := parseRecipient(to)
	if err != nil {
		return permanentError(550, smtp.EnhancedCode{5, 1, 3}, "invalid recipient address")
	}
	if len(s.backend.domains) > 0 {
		if _, ok := s.backend.domains[domain]; !ok {
			return permanentError(550, smtp.EnhancedCode{5, 1, 1}, "recipient domain is not accepted")
		}
	}
	s.recipient = address
	return nil
}

func (s *session) Data(r io.Reader) error {
	if err := s.backend.report(r, s.recipient); errors.Is(err, smtp.ErrDataTooLarge) {
		return smtp.ErrDataTooLarge
	} else if err != nil {
		slog.Error("SMTP delivery failed", "to", s.recipient, "err", err)
		return temporaryError("failed to deliver message")
	}
	return nil
}

func (s *session) Reset() {
	s.recipient = ""
}

func (s *session) Logout() error {
	return nil
}

func (b *backend) report(message io.Reader, recipient string) error {
	reportURL := *b.reportURL
	query := reportURL.Query()
	query.Set("to", recipient)
	reportURL.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodPost, reportURL.String(), message)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	if b.secret != "" {
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		req.Header.Set("X-Tmail-Timestamp", timestamp)
		req.Header.Set("X-Tmail-Signature", sign(b.secret, timestamp, reportURL.Path, recipient))
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	decodeErr := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("report returned HTTP %d: %s", resp.StatusCode, result.Message)
	}
	if decodeErr != nil {
		return fmt.Errorf("decode report response: %w", decodeErr)
	}
	if result.Code != 0 {
		return fmt.Errorf("report failed: %s", result.Message)
	}
	return nil
}

func parseRecipient(value string) (string, string, error) {
	address, err := mail.ParseAddress(value)
	if err != nil {
		return "", "", err
	}
	at := strings.LastIndexByte(address.Address, '@')
	if at <= 0 || at == len(address.Address)-1 {
		return "", "", errors.New("recipient must contain local part and domain")
	}
	domain := normalizeDomain(address.Address[at+1:])
	return address.Address[:at+1] + domain, domain, nil
}

func normalizeDomain(domain string) string {
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)), ".")
}

func sign(secret, timestamp, path, recipient string) string {
	message := timestamp + "\n" + http.MethodPost + "\n" + path + "\n" + recipient
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func temporaryError(message string) error {
	return &smtp.SMTPError{Code: 451, EnhancedCode: smtp.EnhancedCode{4, 3, 0}, Message: message}
}

func permanentError(code int, enhancedCode smtp.EnhancedCode, message string) error {
	return &smtp.SMTPError{Code: code, EnhancedCode: enhancedCode, Message: message}
}
