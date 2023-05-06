AWS_REGION ?= ap-southeast-2
ENVIRONMENT ?= dev
TEMPLATE := template.yml
STACK_NAME := dynamodb-lambda-module-$(ENVIRONMENT)

build:


test:
	cfn-lint $(TEMPLATE) --ignore-checks W3005 W2001

deploy:
	$(eval GITSHA := $(shell git rev-parse --short=7 HEAD ) )
	$(eval AWS_ACCOUNT := $(shell aws sts get-caller-identity --query Account --output text ) )
	$(eval CFN_BUCKET := artificats-s3-$(AWS_ACCOUNT))
	echo Deploying to environment $(ENVIRONMENT)
	aws cloudformation package \
		--template-file $(TEMPLATE) \
		--s3-bucket $(CFN_BUCKET) \
		--s3-prefix $(STACK_NAME) \
		--output-template-file rendered.yml
	aws cloudformation deploy \
		--template-file rendered.yml \
		--stack-name $(STACK_NAME) \
		--s3-bucket $(CFN_BUCKET) \
		--capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM \
		--parameter-overrides Environment=$(ENVIRONMENT) \
			$$( cat params-$(ENVIRONMENT).txt 2>/dev/null ) \
			GitSha=${GITSHA}\
		--no-fail-on-empty-changeset \
		--tags module=dynamodb-lambda-module 

undeploy:
	echo Undeploying from $(ENVIRONMENT)
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME)
	aws cloudformation wait stack-delete-complete \
		--stack-name $(STACK_NAME)

.PHONY: test deploy build
