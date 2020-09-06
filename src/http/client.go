// Copyright (c) 2020, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
  "crypto/tls"
  "fmt"
  "io"
  "io/ioutil"
  "net"
  "net/http"
  "time"
)

type HttpClient struct {
  underlying *http.Client
}

func NewHttpClient() HttpClient {
  return HttpClient{
    underlying: &http.Client{
      Timeout: 120 * time.Second,
      Transport: &http.Transport{
        DialContext: (&net.Dialer{
          Timeout: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout: 10 * time.Second,
        TLSClientConfig: &tls.Config{
          InsecureSkipVerify:       false,
          MinVersion:               tls.VersionTLS12,
          MaxVersion:               tls.VersionTLS12,
          PreferServerCipherSuites: false,
          CurvePreferences: []tls.CurveID{
            tls.CurveP521,
            tls.CurveP384,
            tls.CurveP256,
          },
          CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
            tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
            tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
          },
        },
      },
    },
  }
}

func (client *HttpClient) Get(url string) (response Response, err error) {
  response = Response{
    Status: 0,
    Data:   nil,
    Header: make(map[string]string),
  }

  if client == nil {
    return response, fmt.Errorf("cannot call methods on nil reference")
  }

  var (
    req  *http.Request
    resp *http.Response
  )

  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("runtime error %+v", r)
    }

    if err != nil && resp != nil {
      _, err = io.Copy(ioutil.Discard, resp.Body)
      resp.Body.Close()
    } else if resp == nil && err != nil {
      err = fmt.Errorf("runtime error, no response %+v", err)
    }

    if err == nil {
      response.Data, err = ioutil.ReadAll(resp.Body)
      resp.Body.Close()
    }
  }()

  req, err = http.NewRequest("GET", url, nil)
  if err != nil {
    return
  }
  req.Header.Set("accept", "application/json")
  resp, err = client.underlying.Do(req)
  if err != nil {
    return
  }
  for k, v := range resp.Header {
    response.Header[k] = v[len(v)-1]
  }
  response.Status = resp.StatusCode
  return
}
