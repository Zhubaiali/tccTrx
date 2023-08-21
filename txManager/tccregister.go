package txManager

import "tccTrx/component"

type TCCRegistyCenter interface {
	Register(component component.TCCComponent) error
	Components(componentIDs ...string) ([]component.TCCComponent, error)
}
