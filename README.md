# Pleasant Solutions CLI

Unofficial CLI for [Pleasant products](https://pleasantsolutions.com/), mainly for the [Password server](https://pleasantsolutions.com/passwordserver).


## Pleasant Password server

[Pleasant Password server](https://pleasantsolutions.com/passwordserver) CLI provides a command line integration for the 
password server. Main features:

- Working with folders.
- Working with entries (secrets).


### CLI commands

The base command for the CLI is `pleasant`, specifying Password server as `password-server` or `ps` in short:

```shell
$ pleasant password-server <CMD> <ENTITY> <FLAGS>
$ pleasant ps <CMD> <ENTITY> <FLAGS>
```


### Working with folders

To get information about a Folder:
```shell
# With folder's UUID:
$ pleasant ps get folder <UUID>
# From folder's path:
$ pleasant ps get folder path/to/folder -p
# To print folder's ID: 
$ pleasant ps get folder path/to/folder -p -i
# Check whether a folder exists
$ pleasant ps exists folder path/to/folder
```

To create a new folder:
```shell
$ pleasant ps create folder path/to/folder
```


### Working with entries (secrets)

To get a secret entry:
```shell
# With entry's UUID:
$ pleasant ps get entry <UUID>
# From it's path and name:
$ pleasant ps get entry /path/to/entry/entry-name -f
# To get a web link to an entry:
$ pleasant ps get entry /path/to/entry/entry-name -f -l
```

To create a new secret entry:
```shell
# Create entry 'foo' in the root folder interactively (name of the entry will be user as username as well):
$ pleasant ps create entry foo
# Create entry with random password:
$ pleasant ps create entry foo -r
# Specify username:
$ pleasant ps create entry foo -r -u foo-username
# Place the new entry into specific folder:
$ pleasant ps create entry foo -u foouser -p <folder-UUID>
# Use a path instead of folder UUID - folders will be created if paths doesn't exist:
$ pleasant ps create entry foo -p path/to/folder -f
```

To patch an existing entry:
```shell
# Patch entry identified by UUID with a JSON string:
$ pleasant ps patch entry <UUID> -j '{"Username": "new-username"}'
# Identifying the patching entry with path and name:
$ pleasant ps patch entry path/to/entry/entry-name -f -j '{"Username": "new-username"}'
# Patch with values from a Kubernetes Opaque Secret YAML file:
$ pleasant ps patch entry path/to/entry/entry-name -f -y input.yaml 
```

To clone/duplicate an existing entry, along with all its content and attachments:
```shell
# Clone an entry inside the original folder (adds ' - Copy' to the name):
$ pleasant ps duplicate entry <UUID>
# Clone it to a different folder (identified by the target folder UUID):
$ pleasant ps duplicate entry <UUID> -t <UUID>
# Clone the entry into an existing folder identified by a path:
$ pleasant ps duplicate entry <UUID> -t /foo/bar -p
# Clone the entry into a new folder on given path (i.e. create missing folders in the provided path):
$ pleasant ps duplicate entry <UUID> -t /foo/new-bar -p -c
```

#### Randomized patch values

The [Kubernetes Secret](https://kubernetes.io/docs/concepts/configuration/secret/#opaque-secrets) YAML can contain a 
`$RANDOM$` token as values, e.g.:

```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: fooSecret
stringData:
  randomValue: $RANDOM$
```

The `$RANDOM$` token will be replaced with a random value, if the `randomize` flag is used:

```shell
$ pleasant ps patch entry path/to/entry/entry-name -f -y secret.yaml -r 
```

The `$RANDOM$` token may be even updated in the patch file with actual random value if `update-k8s-opaque-yaml-file` is
used:

```shell
$ pleasant ps patch entry path/to/entry/entry-name -f -y secret.yaml -r -u
```


#### References

The patch file may contain references to existing fields in the Entry. For example if we want to patch a Custom field
with value the `Password`:

```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: fooSecret
stringData:
  referencedValue: $REF.ENTRY.Password.REF$
```


### Global flags

_Global flag might not be implemented for all commands yet!_

- `quiet` modifies output to minimum, so it's easier to parse in scripts.  


### Configuration

Credentials should be places into a `.pleasant.yaml` file, either in `$HOME` or in current folder:

```shell
password_server_url: ""
password_server_username: ""
password_server_password: ""
```
