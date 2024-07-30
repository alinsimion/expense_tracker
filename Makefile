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