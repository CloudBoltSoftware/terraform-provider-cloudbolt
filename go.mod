module github.com/cloudboltsoftware/terraform-provider-cloudbolt

go 1.12

require (
	github.com/cloudboltsoftware/cloudbolt-go-sdk v0.0.0-20191007165942-037df12e87ba
	github.com/hashicorp/terraform v0.12.10
)

// Uncomment the following line if you would like to do local development of the `cloudbolt-go-sdk` library
// Change ../cloudbolt-go-sdk to the path of that local repo
// replace github.com/cloudboltsoftware/cloudbolt-go-sdk => ../cloudbolt-go-sdk
