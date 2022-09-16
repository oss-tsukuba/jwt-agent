Name: jwt-agent
Version: 1.0
Release: 1%{?dist}
Summary: Agent Program for jwt-server
License: BSD
Source0: %{name}-%{version}.tar.gz
BuildRequires: golang

Provides: %{name} = %{version}


%description
Agent for JWT-SERVER

%global debug_package %{nil}

%prep
%setup -q

%build
go mod download golang.org/x/crypto
go get golang.org/x/crypto/ssh/terminal@v0.0.0-20220722155217-630584e8d5aa
go get github.com/mattn/go-isatty
go build jwt-agent.go

%install
rm -rf $RPM_BUILD_ROOT
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

%clean
rm -rf $RPM_BUILD_ROOT

%files
%{_bindir}/%{name}


%changelog
* Mon Aug 29 2022 Atsushi Kumazaki <kuma@canaly.co.jp> 1.0-1
- Initial build.
