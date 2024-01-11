Name: jwt-agent
Version: 1.0.1
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
make

%install
rm -rf $RPM_BUILD_ROOT
make DESTDIR=${RPM_BUILD_ROOT} BINDIR=%{_bindir} MANDIR=%{_mandir} install

%clean
rm -rf $RPM_BUILD_ROOT

%files
%{_bindir}/%{name}
%{_bindir}/jwt-agent-core
%{_mandir}/man1/jwt-agent.1*

%changelog
* Mon Aug 29 2022 Atsushi Kumazaki <kuma@canaly.co.jp> 1.0.0-1
- Initial build.
