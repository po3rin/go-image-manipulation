PHONY: run
run:
	docker run --rm -v ${PWD}:/go/src/gocv-playground \
	gocv-playground /bin/sh -c "${CMD}"
