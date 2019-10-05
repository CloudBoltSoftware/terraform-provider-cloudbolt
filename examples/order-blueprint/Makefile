.PHONY: run

run: init plan apply

init:
	terraform init

plan: init
	terraform plan

destroy: init
	terraform destroy -auto-approve

apply: init
	terraform apply -auto-approve
