// Code generated by "callbackgen -type NRR"; DO NOT EDIT.

package irr

import ()

func (inc *NRR) OnUpdate(cb func(value float64)) {
	inc.updateCallbacks = append(inc.updateCallbacks, cb)
}

func (inc *NRR) EmitUpdate(value float64) {
	for _, cb := range inc.updateCallbacks {
		cb(value)
	}
}
