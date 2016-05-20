**Pivotal Cloud Foundry credential manager CLI helps you configure and interact with deployed credential managers.**

*Starred commands and parameters are planned, but not yet implemented*

```
Usage: cm [<options>] <command> [<args>]
		-v, --version						Show version of CLI and CM API
		-h, --help							Displays help menu


GETTING STARTED: 

 	api -s <server URI>
		View or set the targeted credential manager api
		-s, --server URI					Sets URI for server
		--skip-ssl-validation				Skip verification of the API endpoint. Not recommended!*

	login*
		Authenticates interactively with credential manager.
		-s, --server URL					Sets API target*
		-u, --user USER           			Sets username*
    	-p, --password PASSWORD				Sets password*
		
	logout*
		Discard authenticated user session.


CREDENTIAL MANAGEMENT:

	set --type <cred type> --name <cred name> --secret <cred value>
		Set the value and attributes of a credential.
		-t, --type							Sets the type of credential to store or generate. Default: 'value'*
		-n, --name							Selects the credential being set
		-g, --generate        				System will generate random credential. Cannot be used in combination with --secret.*

		Parameters for setting credential type 'value'
		-s, --secret 						Sets a value for a credential.

		-l, --length NUMBER					Sets length of generated value (Default: 20)*
		--iu, --include-upper FALSE			Sets whether to include UPPER alpha characters (Default: TRUE)*
		--il,  --include-lower FALSE		Sets whether to include lower alpha characters (Default: TRUE)*
		--in, --include-number FALSE		Sets whether to include numeric characters (Default: TRUE)*
		--is, --include-special FALSE		Sets whether to include special characters (Default: TRUE)*

		Parameters for setting credential type 'certificate'
		--ca								Sets the CA value of a certificate credential
		--public-key						Sets the public key value of a certificate credential
		--private-key						Sets the private key value of a certificate credential
		-o, --out							Sets the location to output the generated pem file*
		
		--sca, --signing-ca					Sets the CA used to sign the generated certificate*
		--duration							Sets the valid duration for the generated certificate*
		--kl, --key-length					Sets the bit length of the key*
		--cn, --common-name					Sets the common name of the generated certificate*
		--og, --organization				Sets the organization of the generated certificate*
		--ou, --organization-unit			Sets the organization unit of the generated certificate*
		--lo, --locality					Sets the locality/city of the generated certificate*
		--st, --state						Sets the state/province of the generated certificate*
		--co, --country						Sets the country of the generated certificate*

	get --name <cred name>
		Get the value and attributes of a credential. 
		-n, --name							Selects the credential to retrieve

	delete --name <cred name>
		Delete a credential. 
		-n, --name							Selects the credential to delete
```
