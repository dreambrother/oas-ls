CARGO ?= cargo

.PHONY: check build clean

check:
	$(CARGO) check --target wasm32-wasip2

build:
	$(CARGO) build --release --target wasm32-wasip2

clean:
	$(CARGO) clean
