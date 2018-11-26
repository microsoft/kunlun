# Ansible Role: wildFly

[![Build Status](https://travis-ci.org/bjoernalbers/ansible-role-wildfly.svg?branch=master)](https://travis-ci.org/bjoernalbers/ansible-role-wildfly)

Install [WildFly Application Server](http://wildfly.org) along with OpenJDK on
Ubuntu / openSUSE.

I've developed and tested with Ansible 2.2 on Ubuntu 16.04 and openSUSE 13.
Other versions should work as well (please let me know so that I can update the
supported platforms!)


## Requirements

Your host must have internet access in order to download the WildFly archive.


## Role Variables

### Installed version

Define which version to install (default: WildFly 11):

```yaml
wildfly_version: "10.1.0.Final"
wildfly_checksum: "sha1:5ea0a70a483a4beaf327faeaf0a391208d33d4bd"
```

If you want to install a different version don't forget to update the SHA-1
checksum as well!

### Configuration

You can overwrite these settings:

```yaml
# The configuration you want to run
wildfly_config: standalone.xml

# The mode you want to run
wildfly_mode: standalone

# The address to bind to
wildfly_bind: 0.0.0.0
```

### Wildfly HOME

Please use `wildfly_home` to access Wildfly's basedir, i.e.
`/opt/wildfly-X.Y.Z`.

### Reinstall

Set `wildfly_reinstall` to *yes* when you want to delete `wildfly_home` before
installation.

## Dependencies

None.


## Example Playbook

```yaml
- hosts: servers
  roles:
     - role: bjoernalbers.wildfly
       wildfly_config: standalone-full.xml
```


## License

This Ansible role is released under the [MIT License](LICENSE.txt).
