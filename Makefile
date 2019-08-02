APP="memi"

build:
	@echo "===> Building $(APP)"
	@CGO_ENABLED="0" go build -o ./bin/$(APP)
