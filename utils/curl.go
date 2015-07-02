package pulse

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type CurlResult struct {
	Status    int         //HTTP status of result
	Header    http.Header //Headers
	Remote    string      //Remote IP the connection was made to
	Err       string      //Any Errors that happened. Usually for DNS fail or connection errors.
	Proto     string      //Response protocol
	StatusStr string      //Status in stringified form
}

type CurlRequest struct {
	Path     string
	Endpoint string
	Host     string
	Ssl      bool
}

type MyDialer struct {
	RemoteStr string
}

func (md *MyDialer) Dial(network, address string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	con, err := dialer.Dial(network, address)
	if err == nil {
		md.RemoteStr = con.RemoteAddr().String()
		a, _ := con.RemoteAddr().(*net.TCPAddr)

		if islocalip(a.IP) {
			fmt.Println(a.IP)
			con.Close()
			return nil, securityerr
		}

	}
	return con, err
}

func CurlImpl(r *CurlRequest) *CurlResult {
	var url string
	if r.Ssl {
		url = fmt.Sprintf("https://%s%s", r.Endpoint, r.Path)
	} else {
		url = fmt.Sprintf("http://%s%s", r.Endpoint, r.Path)
	}
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &CurlResult{0, nil, "", err.Error(), "", ""}
	}
	tlshost := r.Endpoint //Validate with endpoint if no host given
	if r.Host != "" {
		req.Host = r.Host
		tlshost = r.Host //Validate with Host hdr if present
	}
	myDialer := &MyDialer{}
	MyTransport := &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: true,
		Dial:              myDialer.Dial,
		ResponseHeaderTimeout: 30 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS10, //TLS 1.0 minimum. Depricating SSLv3 RFC 7568
			ServerName: tlshost,
		}, //Override the hostname to validate
	}
	resp, err := MyTransport.RoundTrip(req)
	if err != nil {
		return &CurlResult{0, nil, myDialer.RemoteStr, err.Error(), "", ""}
	}
	//log.Println(myDialer.RemoteStr)
	//t, _ := http.DefaultTransport.(*http.Transport)
	resp.Body.Close()
	return &CurlResult{resp.StatusCode, resp.Header, myDialer.RemoteStr, "", resp.Proto, resp.Status}
}
