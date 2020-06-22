package msginengine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRsaDecrypt(t *testing.T) {
	smsg := `0lxzMkzb4aIwBv6z21ubUX5UEHnyqfwISiiZeosA4moZrFQ/r5DRVtoGgm2lPmLpiwcrU3H+2Hv+
KOg+iInNId2xGb7qGNKrom071coCIRBu59/wblX2wsWQBVH3MaGk748dSXuqIgayAbAVChEmqMB8
4+aXfLW3Q1mh6nodW7mRxZqNTdhaszWe4IHXSh4H7eIXWZm49bQtf/mi+kQKoEuFAnm3PRbDEeLl
hXv3uF2jRdKewsuWg4Dl7n/r75Mhbu8cfhqMmycZ3zAIjOM4ESNkhKx1u5fMNqk0ruVZGz+eRPJj
z43xd+BPkDmgZg6bNfFn0/ndA4fwlgg5rgsTAg==
`
//	smsg := `q3o1UwIH2MAcYzVV1763UGujsuK9pPMFeqOb/qhh6W+9W6PuDvjmEGKDYaldI7oaxndes4E/jA3cgp9ZZD9CXOEcKJG+MppnC+5qNZIdvHAuM+HLJtmr3FPtdD1yIklF9kKmmzJAyjTvmCtJNhfnPHq4aR4NsIpWKPAa4OaCvWxPKPNQEklJWSAS3OPbKA76FyPbnRudjga9he7t0d7YOduoZqu6+2egyrEdbYvWS4b9RdmcmhE1PfS5PeVb50V5zpWAKqihhSyvVO5KlKgXRrry6sDVMzRvQCXpa8Vqz8GWg3FCOjKk1OzHezA9xswpWWrUa1XQnX3VuHosfKuqSA==
//`
	msg, err := RsaDecrypt(smsg)
	assert.NoError(t, err)

	//fmt.Println(string(msg))
	assert.Equal(t, "this is unsafe", string(msg))
}
