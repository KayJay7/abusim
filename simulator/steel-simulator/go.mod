module steel-simulator

go 1.16

replace (
	steel-lang => ../../src
	steel-simulator-config => ../steel-simulator-config
)
require (
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/containerd/containerd v1.5.4 // indirect
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210722135532-667f2b7c528f // indirect
	gopkg.in/yaml.v2 v2.4.0
	steel-lang v0.0.0
	steel-simulator-config v0.0.0
)
