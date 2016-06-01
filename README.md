**Pivotal Cloud Foundry CredHub CLI helps you configure and interact with deployed CredHub APIs.**

*Starred commands and parameters are planned, but not yet implemented*

```
Usage: cm [<options>] <command> [<args>]
		-v, --version						Show version of CLI and API
		-h, --help							Displays help menu


GETTING STARTED: 

 	api -s <server URI>
		View or set the targeted CredHub API
		-s, --server URI					Sets URI for API server
		--skip-ssl-validation				Skip verification of the API endpoint. Not recommended!*

	login*
		Authenticates interactively with CredHub.
		-s, --server URL					Sets API target*
		-u, --user USER           			Sets username*
    	-p, --password PASSWORD				Sets password*
		
	logout*
		Discard authenticated user session.


CREDENTIAL MANAGEMENT:

	set --type <cred type> --name <cred name> [set params]
		Set the value and attributes of a credential.
		-t, --type							Sets the type of credential to store or generate. Default: 'value'
		-n, --name							Selects the credential being set

		Parameters for setting credential type 'value'
		-v, --value 						Sets the value for the credential.
		
		Parameters for setting credential type 'certificate'
		--ca								Sets the CA based on an input file*
		--public							Sets the public key based on an input file*
		--private							Sets the private key based on an input file*
		--ca-string							Sets the CA to the parameter value*
		--public-string						Sets the public key to the parameter value*
		--private-string					Sets the private key to the parameter value*

	generate --type <cred type> --name <cred name> [generate params]
		Generate and set a credential value based on generation parameters.
		-t, --type							Sets the type of credential to store or generate. Default: 'value'*
		-n, --name							Selects the credential being set*

		Parameters for generating credential type 'value'
		-l, --length NUMBER					Sets length of generated value (Default: 20)*
		--exclude-upper 			        Exclude upper alpha characters from generated value*
		--exclude-lower 		            Exclude lower alpha characters from generated value*
		--exclude-number 		            Exclude numbers from generated value*
		--exclude-special 	  	            Exclude special characters from generated value*

		Parameters for generating credential type 'certificate'
		--signing-ca					Sets the CA used to sign the generated certificate*
		--duration							Sets the valid duration for the generated certificate*
		--key-length					Sets the bit length of the key*
		--common-name					Sets the common name of the generated certificate*
		--organization				Sets the organization of the generated certificate*
		--organization-unit			Sets the organization unit of the generated certificate*
		--locality					Sets the locality/city of the generated certificate*
		--state						Sets the state/province of the generated certificate*
		--country						Sets the country of the generated certificate*
		-o, --out							Sets the location to output the generated pem file*
		
		
	get --name <cred name>
		Get the value and attributes of a credential. 
		-n, --name							Selects the credential to retrieve

	delete --name <cred name>
		Delete a credential. 
		-n, --name							Selects the credential to delete
```
