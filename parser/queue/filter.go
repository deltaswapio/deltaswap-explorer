package queue

import "github.com/deltaswapio/deltaswap/sdk/vaa"

// PythFilter filter vaa event from pyth chain.
func PythFilter(vaaEvent *VaaEvent) bool {
	return vaaEvent.ChainID == uint16(vaa.ChainIDPythNet)
}

// NonFilter non filter vaa evant.
func NonFilter(vaaEvent *VaaEvent) bool {
	return false
}
