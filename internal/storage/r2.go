package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"myplants/internal/config"
)

// r2 Cloudflare R2 兼容 S3 API,这里手写 SigV4 签名,避免引入整个 AWS SDK
type r2 struct {
	cfg    config.R2Config
	host   string
	client *http.Client
}

func NewR2(cfg config.R2Config) *r2 {
	return &r2{
		cfg:    cfg,
		host:   fmt.Sprintf("%s.r2.cloudflarestorage.com", cfg.AccountID),
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *r2) Driver() string { return "r2" }

func (s *r2) Put(key string, data []byte, contentType string) error {
	if contentType == "" {
		contentType = ContentType(key)
	}
	return s.do("PUT", key, data, contentType)
}

func (s *r2) Delete(key string) error {
	return s.do("DELETE", key, nil, "")
}

func (s *r2) do(method, key string, body []byte, contentType string) error {
	path := "/" + s.cfg.Bucket + "/" + encodeKeyPath(key)
	payloadHash := hashHex(body)
	amzDate := time.Now().UTC().Format("20060102T150405Z")
	dateStamp := amzDate[:8]

	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", s.host, payloadHash, amzDate)
	if contentType != "" {
		signedHeaders = "content-type;" + signedHeaders
		canonicalHeaders = "content-type:" + contentType + "\n" + canonicalHeaders
	}

	canonicalRequest := strings.Join([]string{method, path, "", canonicalHeaders, signedHeaders, payloadHash}, "\n")
	scope := dateStamp + "/auto/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256", amzDate, scope, hashHex([]byte(canonicalRequest)),
	}, "\n")

	signature := hex.EncodeToString(hmacBytes(s.signingKey(dateStamp), []byte(stringToSign)))

	req, err := http.NewRequest(method, "https://"+s.host+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s,SignedHeaders=%s,Signature=%s",
		s.cfg.AccessKeyID, scope, signedHeaders, signature))

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		if resp.StatusCode == http.StatusNotFound && method == "DELETE" {
			return nil
		}
		return fmt.Errorf("R2 %s %s: %s %s", method, key, resp.Status, strings.TrimSpace(string(msg)))
	}
	return nil
}

func (s *r2) signingKey(dateStamp string) []byte {
	k := hmacBytes([]byte("AWS4"+s.cfg.SecretAccessKey), []byte(dateStamp))
	k = hmacBytes(k, []byte("auto"))
	k = hmacBytes(k, []byte("s3"))
	return hmacBytes(k, []byte("aws4_request"))
}

func hmacBytes(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func hashHex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// encodeKeyPath 按 S3 规范逐段转义 object key,保留 '/'
func encodeKeyPath(key string) string {
	segs := strings.Split(strings.TrimPrefix(key, "/"), "/")
	for i, s := range segs {
		segs[i] = escapeSeg(s)
	}
	return strings.Join(segs, "/")
}

func escapeSeg(s string) string {
	var b strings.Builder
	for _, c := range []byte(s) {
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' ||
			c == '-' || c == '.' || c == '_' || c == '~' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}
