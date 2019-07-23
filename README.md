# terraform-provider-cloudbolt
Sample Terraform Resource Provider to initiate CloudBolt Blueprint Orders.

## Prerequisites
- Install and Configure golang
- Install and Configure terraform

## Installation
```go
go get github.com/laltomar/cloudbolt-go-sdk
go get github.com/laltomar/terraform-provider-cloudbolt

cd $GO ${GOPATH}/src/github.com/laltomar/terraform-provider-cloudbolt 

mkdir ~/.terraform.d/plugins

go build -o terraform-provider-cloudbolt
mv terraform-provider-cloudbolt ~/.terraform.d/plugins/.

```
