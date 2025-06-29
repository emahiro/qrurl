.PHONY: run
run:
	env $$(cat .env.yaml | grep -v '^#' | grep -v '^$$' | sed 's/: /=/' | xargs) go run main.go