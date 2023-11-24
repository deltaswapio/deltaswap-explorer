package address

import "github.com/deltaswapio/deltaswap-explorer/api/handlers/vaa"

type AddressOverview struct {
	Vaas []*vaa.VaaDoc `json:"vaas"`
}
