NAME="mud client"

ifdef name
NAME=$(name)
endif

install:
	git clone https://github.com/pizzalord22/pizza-client.git
	cd pizza-Client
	go get ./...

compile:
    go build -o $(NAME)


