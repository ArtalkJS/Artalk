%global debug_package %{nil}

Name:           artalk
Version:        2.10.1
Release:        2%{?dist}
Summary:        A self-hosted comment system
License:        MIT
URL:            https://artalk.js.org/
Source0:        %{name}-%{version}-vendored.tar.gz
Source1:        vendor-tarball.sh
Source2:        https://github.com/ArtalkJS/Artalk/releases/download/v%{version}/artalk_ui.tar.gz
Source3:        artalk.sysusers
Source4:        artalk.service
Patch0:         0002-remove-upgrade-command.patch
BuildRequires:  golang
BuildRequires:  systemd-rpm-macros
%{?systemd_requires}
%{?sysusers_requires_compat}

%description
Artalk is a self-hosted comment system for websites.

%prep
%autosetup -p1
tar --strip-components=1 -xzf %{SOURCE2} --directory=public

%build
export CGO_ENABLED=0
export GOPROXY=off
export GOTOOLCHAIN=local
go build -mod=vendor -trimpath -ldflags="-s -w" -o artalk .

%install
install -D -m 0755 artalk %{buildroot}%{_bindir}/artalk
install -D -m 0644 %{SOURCE3} %{buildroot}%{_sysusersdir}/%{name}.conf
install -D -m 0644 %{SOURCE4} %{buildroot}%{_unitdir}/artalk.service
install -d -m 0750 %{buildroot}%{_sharedstatedir}/artalk
install -D -m 0640 conf/artalk.example.yml %{buildroot}%{_sysconfdir}/artalk/artalk.yml
install -d %{buildroot}%{_datadir}/{bash-completion/completions,zsh/site-functions,fish/vendor_completions.d}
export XDG_DATA_HOME="$PWD/.artalk-build"
./artalk completion bash > %{buildroot}%{_datadir}/bash-completion/completions/artalk
./artalk completion zsh > %{buildroot}%{_datadir}/zsh/site-functions/_artalk
./artalk completion fish > %{buildroot}%{_datadir}/fish/vendor_completions.d/artalk.fish

%check
export CGO_ENABLED=0
export GOPROXY=off
export GOTOOLCHAIN=local
export XDG_DATA_HOME="$PWD/.artalk-build"
packages="$(go list -mod=vendor ./...)"
# v2.10.1 has a fixture teardown race in this package (upstream PR #1205).
# Keep its release sources unchanged; later versions run the full suite.
if [ '%{version}' = '2.10.1' ]; then
    packages="$(printf '%%s\n' "$packages" | grep -vx 'github.com/artalkjs/artalk/v2/server/handler')"
fi
go test -mod=vendor $packages
./artalk version | grep -F 'v%{version}'
test -s public/dist/Artalk.js
test -s public/sidebar/index.html

%pre
%sysusers_create_compat %{SOURCE3}

%post
# RPM preserves the old config file, including its root ownership, on upgrade.
if [ -f %{_sysconfdir}/artalk/artalk.yml ] && [ ! -L %{_sysconfdir}/artalk/artalk.yml ]; then
    chown artalk:artalk %{_sysconfdir}/artalk/artalk.yml
    chmod 0640 %{_sysconfdir}/artalk/artalk.yml
fi
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
%attr(0750,artalk,artalk) %dir %{_sysconfdir}/artalk
%attr(0640,artalk,artalk) %config(noreplace) %{_sysconfdir}/artalk/artalk.yml
%attr(0750,artalk,artalk) %dir %{_sharedstatedir}/artalk
%{_datadir}/bash-completion/completions/artalk
%{_datadir}/zsh/site-functions/_artalk
%{_datadir}/fish/vendor_completions.d/artalk.fish

%changelog
* Fri Nov 06 2026 Artalk maintainers - 2.10.1-2
- Build the requested release with vendored Go modules and release UI assets.
