[![CI](https://github.com/heathcliff26/pinentry-keyring/actions/workflows/ci.yaml/badge.svg?event=push)](https://github.com/heathcliff26/pinentry-keyring/actions/workflows/ci.yaml)
[![Editorconfig Check](https://github.com/heathcliff26/pinentry-keyring/actions/workflows/editorconfig-check.yaml/badge.svg?event=push)](https://github.com/heathcliff26/pinentry-keyring/actions/workflows/editorconfig-check.yaml)
[![Renovate](https://github.com/heathcliff26/pinentry-keyring/actions/workflows/renovate.yaml/badge.svg)](https://github.com/heathcliff26/pinentry-keyring/actions/workflows/renovate.yaml)

# pinentry-keyring

Simple wrapper around pinentry to force gpg-agent to display the save passphrase checkbox.

Please note that i wrote this specifically for my use case, so it might not work for you.

## Installing

### From binary

Download the latest release binary for your architecture and mark it as executable.

### Fedora Copr

The app is available as an rpm by using the fedora copr repository [heathcliff26/cli](https://copr.fedorainfracloud.org/coprs/heathcliff26/cli/).
1. Enable the copr repository
```bash
sudo dnf copr enable heathcliff26/cli
```
2. Install the app
```bash
sudo dnf install pinentry-keyring
```

## Configure gpg-agent

Edit the `gpg-agent.conf` and point to the pinenty-keyring binary:
```
vim ~/.gnupg/gpg-agent.conf
```
Add the line:
```
pinentry-program /usr/bin/pinentry-keyring
```

## How does it work

This wrapper ensures that
```
OPTION allow-external-password-cache
```
is set and additionally transforms
```
SETKEYINFO --clear -> SETKEYINFO pinenty-keyring-default-key
```
to ensure the key can be saved.
