// Code generated by "callbackgen -type CA"; DO NOT EDIT.

package indicator

import ()

func (inc *CA) OnUpdate(cb func(value float64)) {
	inc.UpdateCallbacks = append(inc.UpdateCallbacks, cb)
}

func (inc *CA) EmitUpdate(value float64) {
	for _, cb := range inc.UpdateCallbacks {
		cb(value)
	}
}
