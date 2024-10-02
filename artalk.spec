%global goipath github.com/artalkjs/artalk/v2
Version:        2.9.1
%gometa -f
Name:           artalk
Release:        %autorelease
Summary:        A Self-hosted Comment System
License:        MIT
URL:            https://artalk.js.org/
Source0:        %{name}-%{version}-vendored.tar.gz
# track this script
Source1:        vendor-tarball.sh
Source2:        https://github.com/ArtalkJS/Artalk/releases/download/v%{version}/artalk_ui.tar.gz
Source3:        artalk.sysusers
Source4:        artalk.service
BuildRequires:  systemd-rpm-macros
%{?systemd_requires}
%{?sysusers_requires_compat}


%description
Artalk is an intuitive yet feature-rich comment system, ready for immediate deployment into any blog, website, or web application.


%gopkg


%prep
%goprep -k
tar --strip-components=1 -xzf %{SOURCE2} --directory=public
chmod -Rf a+rX,u+w,g-w,o-w .


%build
%undefine _auto_set_build_flags
%global _dwz_low_mem_die_limit 0
%{?gobuilddir:GOPATH="%{gobuilddir}:${GOPATH:+${GOPATH}:}%{?gopath}"} GO111MODULE=on\\
go build %{gobuildflags} -mod=vendor -o %{gobuilddir}/bin/artalk %{goipath}


%install
%gopkginstall

# command
install -m 0755 -vd                     %{buildroot}%{_bindir}
install -m 0755 -vp %{gobuilddir}/bin/* %{buildroot}%{_bindir}/

# sysusers
install -p -D -m 0644 %{SOURCE3} %{buildroot}%{_sysusersdir}/%{name}.conf

# data directory (work dir)
install -d -m 0750 %{buildroot}%{_sharedstatedir}/artalk

# config
install -D -p -m 0644 ./conf/artalk.example.yml %{buildroot}%{_sysconfdir}/artalk/artalk.yml

# systemd units
install -D -p -m 0644 %{SOURCE4} %{buildroot}%{_unitdir}/artalk.service


%pre
%sysusers_create_compat %{SOURCE3}


%post
%systemd_post artalk.service


%preun
%systemd_preun artalk.service


%postun
%systemd_postun_with_restart artalk.service


%files
%doc CONTRIBUTING.md README.md
%license LICENSE
%{_bindir}/artalk
%{_unitdir}/artalk.service
%{_sysusersdir}/%{name}.conf
%dir %{_sysconfdir}/artalk
%config(noreplace) %{_sysconfdir}/artalk/artalk.yml
%attr(0750,artalk,artalk) %dir %{_sharedstatedir}/artalk


%gopkgfiles


%changelog
%autochangelog
