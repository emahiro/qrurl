.PHONY: harness-eval harness-eval-baseline

harness-eval:
	@bash ./harness-eval/run.sh

harness-eval-baseline:
	@bash ./harness-eval/run.sh --write-baseline
