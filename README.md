**Pivotal Cloud Foundry credential manager CLI helps you configure and interact with deployed credential managers.**

*Starred commands and parameters are planned, but not yet implemented*

```
Usage: cm [<options>] <command> [<args>]
		--version				Show version of CLI
		-h, --help				Displays help menu


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


SECRET MANAGEMENT:

	set --name <secret name> --secret <secret value>
		Set the value and attributes of a secret.
		-n, --name							Selects the secret being set
		-s, --secret 						Sets a value for a secret name.
		-t, --type							Sets the type of secret to store or generate. Default: 'value'*

		Populate randomly generated credential using the following parameters:
		-g, --generate        				System will generate random credential. Cannot be used in combination with --secret.*

		Type: 'value' - generation parameters
		-l, --length NUMBER					Sets length of generated value (Default: 20)*
		--iu, --include-upper FALSE			Sets whether to include UPPER alpha characters (Default: TRUE)*
		--il,  --include-lower FALSE		Sets whether to include lower alpha characters (Default: TRUE)*
		--in, --include-number FALSE		Sets whether to include numeric characters (Default: TRUE)*
		--is, --include-special FALSE		Sets whether to include special characters (Default: TRUE)*

		Type: 'certificate' - generation parameters
		-o, --out							Sets the location to output the generated pem file*
		--ca								Sets the CA used to sign the generated certificate*
		--duration							Sets the valid duration for the generated certificate*
		--kl, --key-length					Sets the bit length of the key*
		--cn, --common-name					Sets the common name of the generated certificate*
		--og, --organization				Sets the organization of the generated certificate*
		--ou, --organization-unit			Sets the organization unit of the generated certificate*
		--lo, --locality					Sets the locality/city of the generated certificate*
		--st, --state						Sets the state/province of the generated certificate*
		--co, --country						Sets the country of the generated certificate*

	get --name <secret name>
		Get the value and attributes of a secret. 
		-n, --name							Selects the secret to retrieve

	delete --name <secret name>
		Delete a secret. 
		-n, --name							Selects the secret to delete
```
