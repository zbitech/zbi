ECR_ZBI_REPO = ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
ZBI_REPO=jakinyele
VERSION=1.0.0

build:
	npm run compile

dev:
	npm run dev

image:
	docker build -t ${ZBI_REPO}/zbi-db:latest -t ${ZBI_REPO}/zbi-db:${VERSION} .

image_ecr:
#	export ZBI_REPO=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
#	export VERSION=$(git rev-parse HEAD)
	docker build -t $(ECR_ZBI_REPO)/zbi-db:latest -t $(ECR_ZBI_REPO)/zbi-db:${VERSION} .

push_ecr:
#	export ZBI_REPO=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin $(ECR_ZBI_REPO)
#	export VERSION=$(git rev-parse HEAD)
	docker push $(ECR_ZBI_REPO)/zbi-db:${VERSION}
	docker push $(ECR_ZBI_REPO)/zbi-db:latest
