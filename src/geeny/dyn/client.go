package dyn

import (
	"time"
	"net"
	"net/url"
	"net/http"
	"net/http/cookiejar"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	log "geeny/log"

	"golang.org/x/net/http2"
	"geeny/version"
)

const useHttp2 = true

func NewClient(caPemFile *string, timeout time.Duration, proxy bool) *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	return &http.Client{
		Timeout: timeout,
		Transport: NewTransport(caPemFile, timeout, proxy),
		Jar: jar,
		CheckRedirect: nil,
	}
}

type LoggingTransport struct {
	Transport http.RoundTripper
}

func NewTransport(caPemFile *string, timeout time.Duration, proxy bool) http.RoundTripper {
	var p func(*http.Request) (*url.URL, error) = nil
	if proxy {
		p = http.ProxyFromEnvironment
	}
	delegate := &http.Transport{
		Proxy: p,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).DialContext,
		MaxIdleConns:          1,
		IdleConnTimeout:       timeout,
		TLSHandshakeTimeout:   timeout,
		ExpectContinueTimeout: timeout,
		TLSClientConfig: newTLSClientConfig(caPemFile),
	}
	if useHttp2 {
		err := http2.ConfigureTransport(delegate)
		if err != nil {
			log.Fatal(err)
		}
	}
	return LoggingTransport{delegate}
}

func (t LoggingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add("User-Agent", version.UserAgent)
	req.Header.Add("Connection", "close")

	log.Tracef("> %v %v %v: %v bytes", req.Proto, req.Method, req.URL, req.ContentLength)
	log.Tracef("\tTLS: %+v", req.TLS)
	for k, v := range req.Header {
		log.Tracef("\t%v: %v", k, v)
	}

	resp, err = t.Transport.RoundTrip(req)

	if resp != nil {
		log.Tracef("< %v %v: %v bytes", resp.Proto, resp.Status, resp.ContentLength)
		log.Tracef("\tTLS: %+v", resp.TLS)
		for k, v := range resp.Header {
			log.Tracef("\t%v: %v", k, v)
		}
	}

	return
}

func newTLSClientConfig(caPemFile *string) (cfg *tls.Config) {
	cfg = new(tls.Config)
	if caPemFile != nil {
		cfg.RootCAs = x509.NewCertPool()
		if ca, err := ioutil.ReadFile(*caPemFile); err == nil {
			if ok := cfg.RootCAs.AppendCertsFromPEM(ca); ok {
				log.Tracef("Added CA from %s", caPemFile)
			} else {
				log.Tracef("Adding CA from %s failed", caPemFile)
			}
		} else {
			log.Tracef("Error reading cacert: %v", err)
		}
	}
	return
}

