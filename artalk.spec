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
Source2:        artalk.sysusers
Source3:        https://github.com/ArtalkJS/Artalk/releases/download/v%{version}/artalk_ui.tar.gz
BuildRequires:  systemd-rpm-macros
%{?sysusers_requires_compat}

%description
Artalk is an intuitive yet feature-rich comment system, ready for immediate deployment into any blog, website, or web application.

%gopkg

%prep
%autosetup -c -T -a 1
tar --strip-components=1 -xzf %{SOURCE0}
tar --strip-components=1 -xzf %{SOURCE3} --directory=public
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
install -p -D -m 0644 %{SOURCE2} %{buildroot}%{_sysusersdir}/%{name}.conf
install -d -m 0750 %{buildroot}%{_sharedstatedir}/artalk
install -D -p -m 0644 ./conf/artalk.example.yml %{buildroot}%{_sysconfdir}/artalk/artalk.yml

%pre
%sysusers_create_compat %{SOURCE2}

%files
%doc CONTRIBUTING.md README.md
%license LICENSE
%{_bindir}/artalk
%{_sysusersdir}/%{name}.conf
%dir %{_sysconfdir}/artalk
%config(noreplace) %{_sysconfdir}/artalk/artalk.yml
%attr(0750,artalk,artalk) %dir %{_sharedstatedir}/artalk


%gopkgfiles


%changelog
%autochangelog
