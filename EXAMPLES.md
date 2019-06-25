# CLI Usage Examples

- [CLI Usage Examples](#CLI-Usage-Examples)
  - [Setting](#Setting)
    - [Setting a Value](#Setting-a-Value)
    - [Setting a Password](#Setting-a-Password)
    - [Setting a Certificate](#Setting-a-Certificate)
    - [Setting a JSON blob](#Setting-a-JSON-blob)
    - [Setting a User](#Setting-a-User)
    - [Setting an SSH key](#Setting-an-SSH-key)
    - [Setting an RSA key](#Setting-an-RSA-key)
  - [Getting](#Getting)
    - [Getting the entire value](#Getting-the-entire-value)
    - [Getting only a key](#Getting-only-a-key)
    - [Getting quietly](#Getting-quietly)
    - [Getting a version of a value](#Getting-a-version-of-a-value)
  - [Finding](#Finding)
    - [Find by path](#Find-by-path)
    - [Find by name](#Find-by-name)
  - [Deleting](#Deleting)
    - [Delete a value](#Delete-a-value)
  - [Generating](#Generating)
    - [Generating a Password](#Generating-a-Password)
    - [Generating a Certificate](#Generating-a-Certificate)
    - [Generating a User](#Generating-a-User)
    - [Generating an SSH key](#Generating-an-SSH-key)
    - [Generating an RSA key](#Generating-an-RSA-key)
  - [Permissions](#Permissions)
    - [Setting Permissions](#Setting-Permissions)
    - [Getting Permissions](#Getting-Permissions)
    - [Deleting Permissions](#Deleting-Permissions)
  - [Exporting and Importing](#Exporting-and-Importing)
    - [File Format](#File-Format)
    - [Exporting data](#Exporting-data)
    - [Importing data](#Importing-data)

## Setting

Setting a credential stores the information in CredHub's database for future reference. Used specifically when you know the information you'd like to store ahead of time (i.e. AWS login information).

A global `-t` flag for each set is required and specifies the value type. It can be any of the following: `value`, `password`, `certificate`, `json`, `user`, `ssh`, `rsa`.

A global `-n` flag for each set is also required and specifies the name and path of the value.

### Setting a Value

In addition to the required global flags, the `-v` flag specifies the value you wish to store.

```bash
$ credhub set -t value -n '/example-value' -v 'sample'
id: a3932078-0bc4-41c7-99e5-04af19ae2c32
name: /example-value
type: value
value: <redacted>
version_created_at: "2019-06-24T20:09:43Z"
```

### Setting a Password

In addition to the required global flags, the `-w` flag specifies the password you wish to save.

```bash
$ credhub set -t password -n '/example-password' -w 'supersecret'
id: 6a0471b0-954f-4e22-ade6-e5e1d7e855af
name: /example-password
type: password
value: <redacted>
version_created_at: "2019-06-24T20:16:02Z"
```

### Setting a Certificate

A certificate type in CredHub can store a certificate, private key, and the root CA all in a single entry. Setting a certificate requires at least one of those three items, in addition to the required global flags. The `-c`, `-p`, and `-r` flags (for certificate, private key, and root CA, respectively) can be given a path to a file or a value.

```bash
$ credhub set -t certificate -n '/example-certificate' -c certificate.crt -p private.key -r root_ca.crt
id: caae0e24-d587-44d1-86a5-a3b8649f57d2
name: /example-certificate
type: certificate
value: <redacted>
version_created_at: "2019-06-24T20:19:56Z"
```

If you have a root CA already stored within CredHub, you can reference it when setting a certificate it signed. The `-m` flag should point to the name of the CA stored within CredHub.

```bash
$ credhub set -t certificate -n '/another-certificate' -c certificate.crt -p private.key -m /credhub/path/to/ca
id: 0e662080-b681-4c0c-98c6-057e1c93e771
name: /another-certificate
type: certificate
value: <redacted>
version_created_at: "2019-06-24T20:33:18Z"
```

### Setting a JSON blob

In addition to the required global flags, the `-v` flag specifies the JSON string you wish to store.

```bash
$ credhub set -t json -n '/example-json' -v '{"computer_name": "northport"}'
id: 0fafdac6-4c18-4b7b-84a6-efb1c8f7c072
name: /example-json
type: json
value: <redacted>
version_created_at: "2019-06-25T13:30:12Z"
```

If you wish to read JSON contents from a file, you can use a Bash subshell to set the value:

```bash
$ credhub set -t json -n '/example-json' -v "$(cat my_file.json)"
id: 2c6ca934-5da4-4a68-a617-cb2dcd9e05ef
name: /example-json
type: json
value: <redacted>
version_created_at: "2019-06-25T13:35:51Z"
```

### Setting a User

To set a user, you must use the required global flags and the `-w` flag, which specifies the password. You may also use the `-z` flag to input a username, though it is not required.

```bash
$ credhub set -t user -n '/example-user' -w 'example-password' -z 'example-username'
id: bd49aecd-4d21-4a3b-b28d-3f680c7ca9d2
name: /example-user
type: user
value: <redacted>
version_created_at: "2019-06-25T13:41:07Z"
```

### Setting an SSH key

To set an SSH Key, you must set a name using the `-n` flag, and set one or both of the following: a private key using the `-p` flag, and/or a public key using the `-u` flag. Both private keys and public keys are set from a file or value.

``` bash
$ credhub set -t ssh -n '/example-ssh' -p ./ssh_key.key
id: f07cbd7b-3c31-420f-be4a-31b2267514e1
name: /example-ssh
type: ssh
value: <redacted>
version_created_at: "2019-06-25T13:49:43Z"
```

```bash
$ credhub set -t ssh -n '/example-ssh' -u ./ssh_key.pub
id: f6dfc559-89b5-40cc-ab75-795f5b0e80f4
name: /example-ssh
type: ssh
value: <redacted>
version_created_at: "2019-06-25T13:49:57Z"
```
### Setting an RSA key

To set an RSA Key, you must set a name using the `-n` flag, and set one or both of the following: a private key using the `-p` flag, and/or a public key using the `-u` flag. Both private keys and public keys are set from a file or value.

``` bash
$ credhub set -t rsa -n '/example-rsa' -p ./rsa_key.key
id: f07cbd7b-3c31-420f-be4a-31b2267514e1
name: /example-rsa
type: rsa
value: <redacted>
version_created_at: "2019-06-25T13:49:43Z"
```

```bash
$ credhub set -t rsa -n '/example-rsa' -u ./rsa_key.pub
id: f6dfc559-89b5-40cc-ab75-795f5b0e80f4
name: /example-rsa
type: rsa
value: <redacted>
version_created_at: "2019-06-25T13:49:57Z"
```
## Getting

Getting a credential from CredHub retrieves the information stored in CredHub's database. Usage is the same across value types. 

### Getting the entire value

To get an entire value, input the name of the value using the `-n` flag. The output depends on the type of the value.

```bash
$ credhub get -n '/my-value'
id: 9b3d2e9d-fcbf-4c5c-b15d-644fb263a34f
name: /my-value
type: value
value: hello-world
version_created_at: "2019-06-25T14:02:03Z"
```

As another example, here is the output of getting a certificate value:

```bash
$ credhub get -n '/my-certificate'
id: 4cfa910d-5c49-4c24-a1e7-fc0b65d32922
name: /my-certificate
type: certificate
value:
  ca: |
    -----BEGIN CERTIFICATE-----
    MIIFDzCCAvegAwIBAgIUG9WjrLe06qw5lTxqUxfAwxoteQIwDQYJKoZIhvcNAQEL
    ...
    02tv
    -----END CERTIFICATE-----
  certificate: |
    -----BEGIN CERTIFICATE-----
    MIIFKzCCAxOgAwIBAgIUJaOwSFN805/Y0w9l7b+1Cy7rl4kwDQYJKoZIhvcNAQEL
    ...
    eJCPuoGvHjF/OEOF2zmHqviEp9/s1JQ4QUmJzAS4kQ==
    -----END CERTIFICATE-----
  private_key: |
    -----BEGIN RSA PRIVATE KEY-----
    MIIJKAIBAAKCAgEA2d+5nnqDvve1Ekvz2F7D97ob5u7EJR4Sub2GNh9Tyy9s7ADg
    ...
    A3zYdFe27SW94B+wc15A7nhv5oxRB60559AUpth63DZNTQI1CfMH5XflUJo=
    -----END RSA PRIVATE KEY-----
version_created_at: "2019-06-25T14:06:03Z"
```

### Getting only a key

If you only wish to get a specific key of a stored credential, you can use the `-k` flag, followed by the name of the key. RSA, SSH, Certificate, User, and JSON all have keys you can use. 

```
$ credhub get -n '/my-certificate' -k ca
-----BEGIN CERTIFICATE-----
MIIFDzCCAvegAwIBAgIUG9WjrLe06qw5lTxqUxfAwxoteQIwDQYJKoZIhvcNAQEL
...
6DtQH0VSo06qDZ+weMCemnm8RfepNZfFjGd/Ti4S28FfMZZnQKxPKpcowa0p9O/U
02tv
-----END CERTIFICATE-----
```

In the case of JSON, a value such as this: `{"computer_name": "northport", "boolean": true, "backend": "mysql"}` can be key-retrieved by `computer_name`, `boolean`, and `backend`.

```bash
$ credhub get -n '/my-json' -k 'computer_name'
northport
```
### Getting quietly

The `-q` flag can be put on any get operation to restrict output to only the value.

```bash
$ credhub get -n '/my-value' -q
hello-world
```

### Getting a version of a value

To get more than one version of a value, use the `--versions` flag, and set it equal to the number of past versions you wish to get.

```bash
$ credhub get -n '/test-val' --versions=3
versions:
- id: 60c54e41-3255-4c62-93f3-eb18d8e479c4
  name: /test-val
  type: value
  value: even_newer_val
  version_created_at: "2019-06-25T14:24:04Z"
- id: a3ab3d93-e9f2-4f5b-a12c-0bc3a5058122
  name: /test-val
  type: value
  value: new_val
  version_created_at: "2019-06-25T14:24:00Z"
- id: e2e79a87-c530-482b-8e4f-9a3dcf15cc91
  name: /test-val
  type: value
  value: old_val
  version_created_at: "2019-06-25T14:23:56Z"
```

## Finding

To locate a credential within a CredHub, you can use the `find` command to search for one or more credentials by name or path.

If you wish, you can run `credhub find` to print out the names and version dates of all credentials on the CredHub server. For example:

```bash
$ credhub find
credentials:
- name: /example-json
  version_created_at: "2019-06-25T14:24:04Z"
- name: /example-value
  version_created_at: "2019-06-25T14:14:35Z"
- name: /pcf/password
  version_created_at: "2019-06-25T14:11:04Z"
- name: /personal-certificates/my-certificate
  version_created_at: "2019-06-25T14:06:03Z"
- name: /aws-login
  version_created_at: "2019-06-25T14:05:45Z"
- name: /bosh/bosh-password
  version_created_at: "2019-06-25T14:02:03Z"
- name: /example-ssh
  version_created_at: "2019-06-25T13:49:57Z"
...
```

If the CredHub server is empty, you will receive an empty array: `credentials: []`

### Find by path

To find one or more credentials that start with a specific path, use the `-p` flag followed by the path.

```bash
$ credhub find -p '/example'
credentials:
- name: /example/3
  version_created_at: "2019-06-25T14:44:03Z"
- name: /example/2
  version_created_at: "2019-06-25T14:43:54Z"
- name: /example/1
  version_created_at: "2019-06-25T14:43:41Z"
```

### Find by name

To find one ore more credentials that contain a specific substring, use the `-n` flag followed by the substring.

```bash
$ credhub find -n ample
credentials:
- name: /example/3
  version_created_at: "2019-06-25T14:44:03Z"
- name: /example/2
  version_created_at: "2019-06-25T14:43:54Z"
- name: /example/1
  version_created_at: "2019-06-25T14:43:41Z"
- name: /example-ssh
  version_created_at: "2019-06-25T13:49:57Z"
- name: /example-user
  version_created_at: "2019-06-25T13:44:25Z"
- name: /example-json
  version_created_at: "2019-06-25T13:34:59Z"
- name: /example-certificate
  version_created_at: "2019-06-24T20:19:56Z"
```

## Deleting

Deleting a credential from CredHub removes the value permanently from the database. All versions of the credential are also purged.

### Delete a value

```bash
$ credhub delete -n test-val
Credential successfully deleted
```

## Generating

CredHub can generate credentials if you need a value not previously known (i.e. generating credentials for a new Cloud Foundry platform). They can then be retrieved by the operator using `credhub get` commands. Each value type allows you to set parameters for how the credential should be generated, such as password length or key length.

A global `-t` flag for each generation is required and specifies the value type. It can be any of the following: `value`, `password`, `certificate`, `json`, `user`, `ssh`, `rsa`. 

A global `-n` flag for each generation is also required, and specifies the name and path of the value. 

### Generating a Password

In addition to the required global flags, you can specify password generation parameters, such as length. A table below enumerates the various options you can use. All are optional.

```bash
$ credhub generate -t password -n '/example-password' -l 128 -S
id: 2c99eea4-265b-4bdf-a8b7-d0f8b0b8ccdb
name: /example-password
type: password
value: <redacted>
version_created_at: "2019-06-25T14:53:11Z"
```

Below is a table of the available generation options; please note that some flags are *include* while some are *exclude*. You can use any combination of the flags.

| Flag | Description | Example |
|------|-------------|-------|
|`-l`|Length of the password|`-l 32`
|`-S`|Include special characters in generation| `-S`|
|`-N`|Exclude numbers in generation|`-N`|
|`-U`|Exclude uppercase characters in generation|`-U`|
|`-L`|Exclude lowercase characters in generation|`-L`|

### Generating a Certificate

There's a few groups of required parameters to generate a certificate. The first being the global flags; second being the certificate type and signage, and third being the certificate X509 parameters.

Certificate type and signage can be some of the following. At least one of these flags is required, and only certain combinations are valid (as per X509 spec).

| Flag | Description | Example |
|------|-------------|---------|
|`--ca`|Name of CA in CredHub to sign generated certificate|`--ca /certificate/ca`|
|`--is-ca`|If the certificate should be a CA|`--is-ca`|
|`--self-sign`|If the certificate will be self-signed|`--self-sign`|

X509 parameter flags are below. At least one is required by CredHub to generate the certificate. Any combination of these flags are allowed.

| Flag | Description | Example |
|------|-------------|---------|
|`-k`|Key length of private key|`-k 4096`
|`-d`|Duration of certificate validity|`-d 150`|
|`-c`|Common name of certificate|`-c mydomain.org`|
|`-o`|Organization of certificate|`-o Pivotal`|
|`-u`|Organization unit of certificate|`-u CredHub team`|
|`-i`|Locality of certificate|`-i New York City`|
|`-s`|State of certificate|`-s New York`|
|`-y`|Country of certificate|`-y United States`|
|`-a`|Alternative Name (can be used more than once)|`-a other.mydomain.org`|
|`-g`|Key usage (can be used more than once)|`-g digital_signature`|
|`-e`|Extended key usage (can be used more than once)|`-e encipher_only`|

Below is a table of valid key (and extended key) usage.

|Valid key usage values|
|-|
|`server_auth`|
|`client_auth`|
|`code_signing`|
|`email_protection`|
|`timestamping`|
|`digital_signature`|
|`non_repudiation`|
|`key_encipherment`|
|`data_encipherment`|
|`key_agreement`|
|`key_cert_sign`|
|`crl_sign`|
|`encipher_only`|
|`decipher_only`|

### Generating a User

In addition to the required global flags, you can specify user generation parameters, such as length. A table below enumerates the various options you can use. All are optional.

```bash
$ credhub generate -t user -n '/example-user' -U 
id: 9a0e1b22-fdfa-4961-ba60-a5d884c4a4e4
name: /example-user
type: user
value: <redacted>
version_created_at: "2019-06-25T14:53:11Z"
```

Below is a table of the available generation options; please note that some flags are *include* while some are *exclude*. You can use any combination of the flags.

| Flag | Description | Example |
|------|-------------|-------|
|`-z`|Set a username|`-z example-username`|
|`-l`|Length of the password|`-l 32`
|`-S`|Include special characters in generation| `-S`|
|`-N`|Exclude numbers in generation|`-N`|
|`-U`|Exclude uppercase characters in generation|`-U`|
|`-L`|Exclude lowercase characters in generation|`-L`|



### Generating an SSH key

In addition to the required global flags, there are two optional flags: `-k` and `-m`.

```bash
$ credhub generate -t ssh -n ex-ssh -m gemma@northport -k 4096
id: e9e5ae08-2e1d-498b-a94f-19302388a531
name: /ex-ssh
type: ssh
value: <redacted>
version_created_at: "2019-06-25T16:19:31Z"
```

| Flag | Description | Example |
|------|-------------|---------|
|`-m`| SSH comment|`-m gemma@northport`|
|`-k`|Key length|`-k 4096`|

### Generating an RSA key

In addition to the required global flags, there is one optional flag: `-k`.

```bash
$ credhub generate -t rsa -n ex-rsa -k 3072
id: dcd8b87f-de31-4c02-8b34-6efd72730366
name: /ex-rsa
type: rsa
value: <redacted>
version_created_at: "2019-06-25T16:20:44Z"
```
| Flag | Description | Example |
|------|-------------|---------|
|`-k`|Key length|`-k 3072`|

## Permissions

CredHub supports ACL-based permissioning on credential paths. You can set, delete, and check permissions of paths.

### Setting Permissions

To set permissions of a actor (user) to a specific credential, the `-a` flag specifies the actor you wish to grant permissions to, `-p` specifies the path of the credential, and `o` specifies a list of permissions to grant the actor. Below is a table that enumerates all the available permissions to grant an actor.

```bash
$ credhub set-permission -a garfield -p /example/lasagna -o read,write
actor: garfield
operations:
- read
- write
path: /example/lasagna
uuid: 6a440fb2-c3e5-4469-a8b6-cd40f45bf2ca
```

|Available permissions|Description|
|---------------------|-----------|
|`read`| Permission to read a path/value|
|`write`| Permission to write a path/value|
|`delete`| Permission to delete a path/value|
|`read_acl`| Permission to grant others permission to read a path/value|
|`write_acl`| Permission to grant others permission to write a path/value|


### Getting Permissions

To see if a user has access to a specific credential, the `-a` flag specifies the actor (user) you wish to check and `-p` specifies the path of the credential.

```bash
$ credhub get-permission -a garfield -p /example/lasagna
actor: garfield
operations:
- read
path: /example/lasagna
uuid: 5d320732-8137-4fb6-b002-772a106b5345
```

To see if a user has access to a all values within path, the `-p` flag can alternatively be set to `/path/to/check/*`.

```bash
$ credhub get-permission -a odie -p /example/*
actor: odie
operations:
- read
path: /example/*
uuid: 6a440fb2-c3e5-4469-a8b6-cd40f45bf2ca
```

### Deleting Permissions

Deleting an actor's permissions from a path deletes all operations granted to the user. To do so, the `-a` flag specifies the actor, and the `-p` flag specifies the path or value to remove all permissions granted.

```bash
$ credhub delete-permission -a garfield -p /example/*
actor: garfield
operations:
- read
- write
- delete
- read_acl
- write_acl
path: /example/*
uuid: 6a440fb2-c3e5-4469-a8b6-cd40f45bf2ca
```

## Exporting and Importing

CredHub allows you to bulk export or import values as a YAML file. 

### File Format

```yaml
credentials:
- name: /example-rsa
  type: rsa
  value:
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      ...
      -----END RSA PRIVATE KEY-----
    public_key: |
      -----BEGIN PUBLIC KEY-----
      ...
      -----END PUBLIC KEY-----
- name: /example-ssh
  type: ssh
  value:
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      ...
      -----END RSA PRIVATE KEY-----
    public_key: ssh-rsa ...
    public_key_fingerprint: ...
- name: /example-user
  type: user
  value:
    password: axdm3qmdfoz84l90bye5ae4lfhfflq
    password_hash: $6$HSkXv3EA$eQpxsEc0trr1l21ec7vwXkw.691oP.vn6I//
    username: GgRTWSNzBKXYGTIRErvk
- name: /example-password
  type: password
  value: WVy6bzRNQam0nFuApbbIBHqXCvqqzVWygASZYW6kovFCM8DszAzsDiuD1UT
- name: /example-value
  type: value
  value: "hello"
- name: /final-cert
  type: certificate
  value:
    ca: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    certificate: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      ...
      -----END RSA PRIVATE KEY-----
- name: /json_value
  type: json
  value:
    name: "garfield"
    computer_name: "northport"
    backend: "mysql"
```

### Exporting data

The `-f` flag specifies the output name of the file.

```bash
$ credhub export -f my_credhub_export.yml
```

### Importing data

The `-f` flag specifies the file to read from.

```bash
$ credhub import -f my_credhub_export.yml

Import complete.
Successfully set: 21
Failed to set: 0
```