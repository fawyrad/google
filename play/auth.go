package play

import (
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

type GoogleAuth struct {
   Values Values
}

func (g GoogleAuth) auth() string {
   return g.Values["Auth"]
}

func (g *GoogleToken) token() string {
   return g.Values["Token"]
}

type GoogleToken struct {
   Values Values
}

func (GoogleToken) Marshal(token string) ([]byte, error) {
   resp, err := http.PostForm(
      "https://android.googleapis.com/auth", url.Values{
         "ACCESS_TOKEN": {"1"},
         "Token":        {token},
         "service":      {"ac2dm"},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   return io.ReadAll(resp.Body)
}

func (g *GoogleToken) Unmarshal(data []byte) error {
   g.Values = Values{}
   return g.Values.Set(string(data))
}

func (g *GoogleToken) Auth() (*GoogleAuth, error) {
   resp, err := http.PostForm(
      "https://android.googleapis.com/auth", url.Values{
         "Token":      {g.token()},
         "app":        {"com.android.vending"},
         "client_sig": {"38918a453d07199354f8b19af05ec6562ced5788"},
         "service":    {"oauth2:https://www.googleapis.com/auth/googleplay"},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   query := Values{}
   query.Set(string(data))
   return &GoogleAuth{query}, nil
}

func (v Values) Set(query string) error {
   for query != "" {
      var key string
      key, query, _ = strings.Cut(query, "\n")
      key, value, _ := strings.Cut(key, "=")
      v[key] = value
   }
   return nil
}

type Values map[string]string
