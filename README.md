# CredHub CLI 

CredHub manages credentials like passwords, certificates, ssh keys, rsa keys, strings (arbitrary values) and CAs. CredHub provides a CLI and API to get, set, generate and securely store such credentials.

* [CredHub Tracker](https://www.pivotaltracker.com/n/projects/1977341)

See additional repos for more info:

* [credhub](https://github.com/cloudfoundry-incubator/credhub) :     CredHub server code 
* [credhub-acceptance-tests](https://github.com/cloudfoundry-incubator/credhub-acceptance-tests) : Integration tests
* [credhub-release](https://github.com/pivotal-cf/credhub-release) : BOSH release of CredHub server

### Building the CLI:

`make` (first time only to get dependencies, will also run specs)

`make build`


### Usage:

```
Usage: credhub [<options>] <command> [<args>]
    --version                            Show version of CLI and API
    -h, --help                           Displays help menu


GETTING STARTED: 

  api
    View or set the targeted CredHub API (short command: a)
    -s, --server 'URI'                   URI of API server to target
        --skip-tls-validation            Skip certificate validation of the API endpoint. Not recommended!

  login
    Authenticates interactively with CredHub (short command: l)
    -u, --user 'USER'                    Authentication username
    -p, --password 'PASSWORD'            Authentication password
    -s, --server 'URI'                   URI of API server to target
        --skip-tls-validation            Skip certificate validation of the API endpoint. Not recommended!

  logout
    Discard authenticated user session (short command: o)


CREDENTIAL MANAGEMENT:

  get --name <cred name>
    Get the value and attributes of a Credential (short command: g)
    -n, --name 'CRED'                     Name of the credential to retrieve
        --output-json                     Return response in json format

  set --type <cred type> --name <cred name> [set params]
    Set the value and attributes of a credential. Supported types 'password', 'value', 'certificate', 'ssh' and 'rsa' (short command: s)
    -n, --name 'CRED'                     Name of the credential to set
    -t, --type [TYPE]                     Sets the credential type (Default: 'password')
    -O, --no-overwrite                    Credential is not modified if stored value already exists

    Set parameters by [Type]
    -v, --value 'VALUE'                   [Password, Value] Sets the value for the credential
    -r, --root  <FILE>                    [Certificate] Sets the root CA from file
    -c, --certificate <FILE>              [Certificate] Sets the certificate from file
    -p, --private <FILE>                  [Certificate, SSH] Sets the private key from file
    -u, --public <FILE>                   [SSH, RSA] Sets the public key from file
    -R, --root-string 'ROOT'              [Certificate] Sets the root CA from string input
    -C, --certificate-string 'CERT'       [Certificate] Sets the certificate from string input
    -P, --private-string 'PRIVATE'        [Certificate, SSH] Sets the private key from string input
    -U, --public-string 'PUBLIC'          [SSH, RSA] Sets the public key from string input

  generate --type <cred type> --name <cred name> [generate params]
    Generate and set a credential value based on generation parameters. Supported types 'password', 'certificate', 'ssh' and 'rsa'  (short command: n)
    -n, --name 'CRED'                      Name of the credential to generate
    -t, --type [TYPE]                      Sets the credential type to generate (Default: 'password')
    -O, --no-overwrite                     Credential is not modified if stored value already exists
        --output-json                      Return response in json format

    Generate parameters by [Type]
    -l, --length [4-200]                   [Password] Length of generated value (Default: 30)
    -U, --exclude-upper                    [Password] Exclude upper alpha characters from generated value
    -L, --exclude-lower                    [Password] Exclude lower alpha characters from generated value
    -N, --exclude-number                   [Password] Exclude number characters from generated value
    -S, --include-special                  [Password] Include special characters in the generated value
    -H, --only-hex                         [Password] Use only hexadecimal characters in generated value   
    -k, --key-length [2048, 3072, 4096]    [Certificate, SSH, RSA] Bit length of the generated key (Default: 2048)
    -m, --ssh-comment 'COMMENT'            [SSH] A comment appended to the SSH public key for identification 
    -d, --duration [1-3650]                [Certificate] Valid duration (in days) of the generated certificate (Default: 365)
    -c, --common-name 'COMMON NAME'        [Certificate] Common name of the generated certificate 
    -a, --alternative-name 'ALT NAME'      [Certificate] A subject alternative name of the generated certificate (may be specified multiple times)
    -o, --organization 'ORG'               [Certificate] Organization of the generated certificate
    -u, --organization-unit 'ORG UNIT'     [Certificate] Organization unit of the generated certificate
    -i, --locality 'LOCALITY'              [Certificate] Locality/city of the generated certificate
    -s, --state  'ST'                      [Certificate] State/province of the generated certificate
    -y, --country 'CC'                     [Certificate] Country of the generated certificate
    -g, --key-usage [key_usage]            [Certificate] Key Usage extensions for the generated certificate (may be specified multiple times)
    -e, --ext-key-usage [ext_key_usage]    [Certificate] Extended Key Usage extensions for the generated certificate (may be specified multiple times)
        --ca 'CA NAME'                     [Certificate] Name of CA used to sign the generated certificate
        --self-sign                        [Certificate] The generated certificate will be self-signed
        --is-ca                            [Certificate] The generated certificate is a certificate authority
    
  regenerate --name <cred name>
    Regenerates a credential using the same parameters that were previously used (short command: r)
    -n, --name 'CRED'                      Name of the credential to regenerate

  delete --name <cred name>
    Delete a credential (short command: d)
    -n, --name 'CRED'                      Name of the credential to delete
    
  find [find params]
    Find existing credentials based on query parameters (short command: f)
    -n, --name-like 'CRED'                 Find credentials by partial name search
    -p, --path 'PATH'                      Find credentials by path
    -a, --all-paths                        List all existing credential paths
```
