# vim: set ft=ruby ts=2 sw=2:

options = {
  :name => 'gofr',
  :command => 'make daemon',
  :env => {},
}

interactor :off

guard( :process, options ) do
  watch( /.*\.(go)$/ )
end
