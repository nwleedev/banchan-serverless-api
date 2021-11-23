.PHONY: build clean deploy gomodgen

build:
	go build -ldflags="-s -w" -o bin/products product/list/main.go
	go build -ldflags="-s -w" -o bin/product product/item/main.go
	go build -ldflags="-s -w" -o bin/product_search product/search/main.go
	go build -ldflags="-s -w" -o bin/tags tag/list/main.go
	go build -ldflags="-s -w" -o bin/tag tag/item/main.go
	go build -ldflags="-s -w" -o bin/tag_search tag/search/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy:
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
