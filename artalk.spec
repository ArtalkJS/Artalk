%global goipath github.com/artalkjs/artalk/v2
Version:        2.9.1
%gometa -f
Name:           artalk
Release:        %autorelease
Summary:        A Self-hosted Comment System

License:        MIT
URL:            https://artalk.js.org/
Source0:        %{gosource}
Source1:        vendor-%{version}.tar.gz

%description
Artalk is an intuitive yet feature-rich comment system, ready for immediate deployment into any blog, website, or web application.

%gopkg

%prep
%autosetup -c -T -a 1
tar --strip-components=1 -xzf %{SOURCE0}
chmod -Rf a+rX,u+w,g-w,o-w .
%goprep -k -e


%build
%undefine _auto_set_build_flags
%global _dwz_low_mem_die_limit 0
%{?gobuilddir:GOPATH="%{gobuilddir}:${GOPATH:+${GOPATH}:}%{?gopath}"} GO111MODULE=on\\
go build %{gobuildflags} -mod=vendor -o %{gobuilddir}/bin/artalk %{goipath}


%install
%gopkginstall
install -m 0755 -vd                     %{buildroot}%{_bindir}
install -m 0755 -vp %{gobuilddir}/bin/* %{buildroot}%{_bindir}/


%files
%doc CONTRIBUTING.md README.md
%license LICENSE
%{_bindir}/artalk


%gopkgfiles


%changelog
%autochangelog
