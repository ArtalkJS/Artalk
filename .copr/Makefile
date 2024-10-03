tag=

srpm:
	dnf -y install git rpmdevtools golang rpmautospec rpkg
	$(eval tag=$(shell git tag -l 'v*' --sort=-v:refname | head -n 1))
	@echo $(tag)
	spectool -g artalk.spec
	./vendor-tarball.sh ${tag}
	rpmautospec process-distgit artalk.spec artalk.spec
	rpkg srpm --spec artalk.spec --outdir=$(outdir)
