module github.com/matnich89/benefex/distribution

go 1.21

require (
	github.com/matnich89/benefex/common v0.0.0
	github.com/rabbitmq/amqp091-go v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/matnich89/benefex/common => ../common
