BUILD_VERSION   := v1.0.1
BUILD_TIME      := $(shell date "+%F %T")
# BUILD_NAME      := SmartControlServer_$(shell date "+%Y%m%d%H" )
BUILD_NAME      := easydarwin
SOURCE          := ./
TARGET_DIR      := ./
COMMIT_SHA1     := $(shell git rev-parse HEAD )

default: release


debug:
    # CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -ldflags                   \
	"                                           \
	-X 'main.BuildVersion=${BUILD_VERSION}'     \
	-X 'main.BuildTime=${BUILD_TIME}'       \
	-X 'main.BuildName=${BUILD_NAME}'       \
	-X 'main.CommitID=${COMMIT_SHA1}'       \
	"                                           \
	-o ${BUILD_NAME} ${SOURCE}

release:
    # CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -tags release -ldflags              \
	"                                           \
	-X 'main.BuildVersion=${BUILD_VERSION}'     \
	-X 'main.BuildTime=${BUILD_TIME}'       \
	-X 'main.BuildName=${BUILD_NAME}'       \
	-X 'main.CommitID=${COMMIT_SHA1}'       \
	"                                           \
	-o ${BUILD_NAME} ${SOURCE}

install:
	mkdir -p ${TARGET_DIR}
	cp ${BUILD_NAME} ${TARGET_DIR} -f


clean:
	rm ${BUILD_NAME} -f


.PHONY : all clean install ${BUILD_NAME}