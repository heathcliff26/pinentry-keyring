%global debug_package %{nil}

Name:           pinentry-keyring
Version:        0
Release:        %autorelease
Summary:        Wrapper around pinentry to fix gpg-agent
%global package_id io.github.heathcliff26.%{name}

License:        Apache-2.0
URL:            https://github.com/heathcliff26/%{name}
Source:         %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires: golang >= 1.24

%global _description %{expand:
Lightweight wrapper around pinentry to ensure gpg shows the option to save the passphrase to gnome keyring.}

%description %{_description}

%prep
%autosetup -n %{name}-%{version} -p1

%build
export RELEASE_VERSION="%{version}-%{release}"
make build

%install
install -D -m 0755 bin/%{name} %{buildroot}%{_bindir}/%{name}
install -D -m 0644 %{package_id}.metainfo.xml %{buildroot}/%{_datadir}/metainfo/%{package_id}.metainfo.xml

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}
%{_datadir}/metainfo/%{package_id}.metainfo.xml

%changelog
%autochangelog
