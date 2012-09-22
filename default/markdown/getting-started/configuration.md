# Configuring Luminos

Luminos has a ``settings.yaml`` file that may look like this:


```yaml
# This is a YAML file.
# Please use two-space tabs instead of the \t character.
server:
  bind: "0.0.0.0"       # Address to bind to, listen all by default.
  port: 9000            # Port for listening to.
  type: "standalone"    # Types are "fastcgi" or "standalone".

hosts:
  default: "./default"
# beta.example.org: /home/user/projects/beta.example.org
```

By default, Luminos listens on all interfaces on the 9000 port and all the request that have an unknown
virtual host are looked up into the ``./default`` directory.

To add a virtual host, just add a name under the ``host:`` definition of the ``settings.yaml`` file and
point it to a directory.

Directories must have an [special structure](/getting-started/directory-structure) in order to work with
Luminos.

## Using another settings file

By default, Luminos looks into the current working directory for a ``settings.yaml`` file, but this
can be configured:

```bash
% ./luminosd -conf /etc/luminos/settings.yaml
```
