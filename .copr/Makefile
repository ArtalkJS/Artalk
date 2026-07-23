srpm:
	.copr/copr_script.sh "$(outdir)" "$(version)"
	rpmbuild -bs --define "_sourcedir $(abspath $(outdir))" --define "_srcrpmdir $(abspath $(outdir))" "$(outdir)/artalk.spec"
