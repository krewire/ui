module github.com/krewire/ui/core

go 1.27.3

require (
	github.com/krewire/ui/components v0.0.0
	github.com/krewire/ui/form v0.0.0
	github.com/krewire/ui/panel v0.0.0
	github.com/krewire/ui/widget v0.0.0
	github.com/krewire/ui/kiw v0.0.0
	github.com/krewire/ui/node v0.0.0
	github.com/krewire/ui/dsl v0.0.0
	github.com/krewire/ui/editors v0.0.0
	github.com/krewire/ui/styles v0.0.0
)

replace (
	github.com/krewire/ui/components => ../components
	github.com/krewire/ui/form => ../form
	github.com/krewire/ui/panel => ../panel
	github.com/krewire/ui/widget => ../widget
	github.com/krewire/ui/kiw => ../kiw
	github.com/krewire/ui/node => ../node
	github.com/krewire/ui/dsl => ../dsl
	github.com/krewire/ui/editors => ../editors
	github.com/krewire/ui/styles => ../styles
)
