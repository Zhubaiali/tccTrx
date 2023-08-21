package txManager

import "tccTrx/component"

type TCCRegistryCenter interface {
	Register(component component.TCCComponent) error
	Components(componentIDs ...string) ([]component.TCCComponent, error)
}
