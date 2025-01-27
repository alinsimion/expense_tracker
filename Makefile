build:
	go build -o bin/tracker

generate:
	~/go/bin/templ generate

tailwind:	
	npx tailwindcss -i ./css/style.css -o ./css/output.css

# tailwind 
run: generate build 
	./bin/tracker

# test: 
# 	cd ./types && go test -v

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/tracker main.go

build-to-deploy: tailwind generate build build-linux

deploy-all:
	scp -r "/Users/asimion/Desktop/Personal/Projects/expense_tracker/bin" root@143.42.59.21:/home/expense_tracker
	scp -r "/Users/asimion/Desktop/Personal/Projects/expense_tracker/static" root@143.42.59.21:/home/expense_tracker/bin/
