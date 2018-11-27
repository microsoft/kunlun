Ansible Role: Apache JMeter 3.3
=========
An Ansible Role that install Apache JMeter

The Apache JMeter™ application is open source software, a 100% pure Java application designed to load test functional behavior and measure performance. It was originally designed for testing Web Applications but has since expanded to other test functions.

Requirements
------------
JMeter requires a fully compliant JVM 8, we advise that you install latest minor version of those major versions. Java 9 is not tested completely as of JMeter 3.2.

Role Variables
--------------

A description of the settable variables for this role should go here, including any variables that are in defaults/main.yml, vars/main.yml, and any variables that can/should be set via parameters to the role. Any variables that are read from other roles and/or the global scope (ie. hostvars, group vars, etc.) should be mentioned here as well.

Dependencies
------------

None. Note: JMeter has JAVA dependency. 

Example Playbook
----------------

Including an example of how to use your role (for instance, with variables passed in as parameters) is always nice for users too:

    - hosts: servers
      roles:
         - { role: wikivoks.jmeter }

License
-------

BSD

Author Information
------------------

This role was created by Wikivoks. 
