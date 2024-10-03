src := $(shell pwd)/.copr
tag=

srpm: dep
	pushd $(src)
	$(eval tag=$(shell git tag -l 'v*' --sort=-v:refname | head -n 1))
	@echo $(tag)
	spectool -g artalk.spec
	./vendor-tarball.sh $(tag)
	git stash					# make sure workspace is empty.
	cp changelog artalk.spec ../
	git add ../changelog ../artalk.spec && git commit -m "for rpmautospec [skip changelog]"
	rpmautospec process-distgit ../artalk.spec artalk.spec
	git stash pop				# pop it!
	rpkg srpm --spec artalk.spec --outdir=$(outdir)
	popd

dep:
	dnf -y install git rpmdevtools golang rpmautospec rpkg
