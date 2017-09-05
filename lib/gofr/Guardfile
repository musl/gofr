# vim: set ft=ruby sw=4 ts=4 :

clearing :on
interactor :off

opts = {
	name: 'gofr',
	env: {},
	command: 'make',
}

guard( :process, opts ) do
	watch( '.*\.(go|sh)$' )
	watch( 'Makefile' )
end

