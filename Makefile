OUTPUT_PATH=./
PROGRAMS = suggest_web

all: clean default

test: check

default:
	    go build -i -o $(OUTPUT_PATH)/suggest_web ./suggest_web.go

linux:
	    GOOS=linux go build -i -o $(OUTPUT_PATH)/suggest_web ./suggest_web.go

check:
	    go test client/dataapi*

clean:
	    @rm -rf $(PROGRAMS)
