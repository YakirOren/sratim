package sratim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	TestToken = "bruh"
	TestID    = "6969"
	TestData  = "sometestdata"
)

type SratimTestSuite struct {
	suite.Suite
	client *Sratim

	main *httptest.Server
	api  *httptest.Server
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPluginSuite(t *testing.T) {
	suite.Run(t, new(SratimTestSuite))
}

func (suite *SratimTestSuite) SetupTest() {
	suite.main = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		http.SetCookie(res, &http.Cookie{Name: "bruh", Value: "sdfsldfjsdklfj"})

		_, err := res.Write([]byte(`bruh`))
		if err != nil {
			require.Nil(suite.T(), err)
		}
	}))

	suite.api = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var response []byte

		if req.URL.String() == "/bruh" {
			res.Header().Set("Content-Type", "text/plain")
			response = []byte(TestData)
		}

		if req.URL.String() == "/movie/preWatch" {
			res.Header().Set("Content-Type", "text/plain")
			response = []byte(TestToken)
		}

		if req.URL.String() == "/movie/watch/id/"+TestID+"/token/"+TestToken {
			res.Header().Set("Content-Type", "application/json")
			response, _ = json.Marshal(Response{
				Success: true,
				Watch: struct {
					URL string `json:"480"`
				}{URL: fmt.Sprintf(`//s1.sratim.tv/movie/SD/480/%s.mp4?token=%s&time=1639066886&uid=`, TestID, TestToken)},
			})

		}

		res.WriteHeader(http.StatusOK)
		_, err := res.Write(response)
		if err != nil {
			require.Nil(suite.T(), err)
		}
	}))
	suite.client, _ = New(suite.main.URL, suite.api.URL)
}

func (suite *SratimTestSuite) TeardownTest() {
	suite.main.Close()
	suite.api.Close()
}

func (suite *SratimTestSuite) TestNew() {
	type args struct {
		url    string
		apiUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		token   string
	}{
		{
			name: "create new",
			args: args{
				url:    suite.main.URL,
				apiUrl: suite.api.URL,
			},

			token:   TestToken,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.url, tt.args.apiUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.token, tt.token) {
				t.Errorf("New() got = %v, want %v", got, tt.token)
			}
		})
	}
}

func (suite *SratimTestSuite) TestSratim_GetMovieURL() {
	u, _ := url.Parse(fmt.Sprintf("https://s1.sratim.tv/movie/SD/480/%s.mp4?token=%s&time=1639066886&uid=", TestID, TestToken))

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		client  *Sratim
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "get movie url",
			client: suite.client,
			args: args{
				id: TestID,
			},
			want:    u.String(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			sr := Sratim{
				client: tt.client.client,
				token:  tt.client.token,
				url:    tt.client.url,
				apiUrl: tt.client.apiUrl,
			}
			got, err := sr.GetMovieURL(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMovieURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("%v", err)
				return
			}
			if got.String() != tt.want {
				t.Errorf("GetMovieURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *SratimTestSuite) TestSratim_download() {
	tests := []struct {
		name       string
		fields     *Sratim
		movieURL   string
		wantWriter string
		wantErr    bool
	}{
		{
			name:       "get some data",
			fields:     suite.client,
			movieURL:   suite.api.URL + "/bruh",
			wantWriter: TestData,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			sr := Sratim{
				client: tt.fields.client,
				token:  tt.fields.token,
				url:    tt.fields.url,
				apiUrl: tt.fields.apiUrl,
			}
			writer := &bytes.Buffer{}
			err := sr.download(tt.movieURL, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("download() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
