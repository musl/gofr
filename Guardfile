clearing :off
interactor :off

guard(:process, name: 'gofrd', command: 'make -C cmd/gofrd clean run') do
  watch(/.*\.go$/)
end

guard(:process, name: 'node', command: 'make -C cmd/gofrd dev')

