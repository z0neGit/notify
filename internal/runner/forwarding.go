package runner

import (
	"fmt"
	"github.com/projectdiscovery/gologger"
	"io/ioutil"
	"net/http"
)

type ForwardingHandler struct {
	forwarderChan chan string
	authToken     string
}

func (fh *ForwardingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqAuthToken := r.URL.Query().Get("token")
	if fh.authToken != "" && reqAuthToken != fh.authToken {
		w.WriteHeader(403)
		return
	}

	fmt.Fprintf(w, "OK")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		gologger.Error().Msgf("%s\n", err)
	}
	fh.forwarderChan <- string(bodyBytes)
}

func startForwarding(listen string, forwarderChan chan string, authToken string) {
	http.Handle("/forwarding", &ForwardingHandler{forwarderChan: forwarderChan, authToken: authToken})
	go func() {
		if authToken == "" {
			gologger.Info().Msgf("Start forwarding service with http://%s/forwarding\n", listen)
		} else {
			gologger.Info().Msgf("Start forwarding service with http://%s/forwarding?token=%s\n", listen, authToken)
		}

		err := http.ListenAndServe(listen, nil)
		if err != nil {
			gologger.Fatal().Msgf("%s\n", err)
		}
	}()
}
