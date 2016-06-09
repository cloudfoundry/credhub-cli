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
		-u, --user USER           			Sets username*
    	-p, --password PASSWORD				Sets password*
    	-s, --server URL					Sets API target*

		
	logout*
		Discard authenticated user session.


CREDENTIAL MANAGEMENT:

	set --type <cred type> --name <cred name> [set params]
		Set the value and attributes of a credential.
		-t, --type							Sets the type of credential to store or generate. (Default: 'value')
		-n, --name							Selects the credential being set

		Set parameters by [Type]
		-v, --value 						[Value] Sets the value for the credential.
		--ca								[Certificate] Sets the CA based on an input file*
		--public							[Certificate] Sets the public key based on an input file*
		--private							[Certificate] Sets the private key based on an input file*
		--ca-string							[Certificate] Sets the CA to the parameter value
		--public-string						[Certificate] Sets the public key to the parameter value
		--private-string					[Certificate] Sets the private key to the parameter value

	generate --type <cred type> --name <cred name> [generate params]
		Generate and set a credential value based on generation parameters.
		-t, --type							Sets the type of credential to store or generate. (Default: 'value')
		-n, --name							Selects the credential being set

		Set parameters by [Type]
		-l, --length NUMBER					[Value] Sets length of generated value (Default: 20)
		--exclude-upper 			        [Value] Exclude upper alpha characters from generated value
		--exclude-lower 		            [Value] Exclude lower alpha characters from generated value
		--exclude-number 		            [Value] Exclude numbers from generated value
		--exclude-special 	  	            [Value] Exclude special characters from generated value
		--ca					     	    [Certificate] Sets the CA used to sign the generated certificate*
		--duration							[Certificate] Sets the valid duration for the generated certificate*
		--key-length						[Certificate] Sets the bit length of the key*
		--common-name						[Certificate] Sets the common name of the generated certificate*
		--alternate-name					[Certificate] Sets an alternate name of the generated certificate*
		--organization						[Certificate] Sets the organization of the generated certificate*
		--organization-unit					[Certificate] Sets the organization unit of the generated certificate*
		--locality							[Certificate] Sets the locality/city of the generated certificate*
		--state								[Certificate] Sets the state/province of the generated certificate*
		--country							[Certificate] Sets the country of the generated certificate*
		
		
	get --name <cred name>
		Get the value and attributes of a credential. 
		-n, --name							Selects the credential to retrieve

	delete --name <cred name>
		Delete a credential. 
		-n, --name							Selects the credential to delete
```
