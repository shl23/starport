package cosmosfaucet

import "net/http"

type transferRequest struct {
	AccountAddress string `json:"address"`
}

type transferResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// ServeHTTP implements http.Handler to expose the functionality of Faucet.Transfer() via HTTP.
// request/response payloads are compatible with the previous implementation at allinbits/cosmos-faucet.
func (f Faucet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
