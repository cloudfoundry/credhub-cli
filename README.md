**Pivotal Cloud Foundry CredHub CLI helps you configure and interact with deployed CredHub APIs.**

*Starred commands and parameters are planned, but not yet implemented*

```
Usage: cm [<options>] <command> [<args>]
		-v, --version					 		Show version of CLI and API
		-h, --help								Displays help menu


GETTING STARTED: 

 	api -s <server URI>
		View or set the targeted CredHub API
		-s, --server 'URI'						URI of API server to target

	login*
		Authenticates interactively with CredHub.
		-u, --user 'USER'           			Authentication username*
    	-p, --password 'PASSWORD'				Authentication password*
		-s, --server 'URI'						URI of API server to target*

	logout*
		Discard authenticated user session.


CREDENTIAL MANAGEMENT:

	get --name <cred name>
		Get the value and attributes of a credential. 
		-n, --name 'CRED'						Name of credential to retrieve

	set --type <cred type> --name <cred name> [set params]
		Set the value and attributes of a credential.
		-t, --type ['value', 'certificate']		Sets the type of credential to store or generate. (Default: 'value')
		-n, --name 'CRED'						Selects the credential being set

		Set parameters by [Type]
		-v, --value 'VALUE'						[Value] Sets the value for the credential.
		--root	<FILE>							[Certificate] Sets the root CA from file
		--certificate <FILE>					[Certificate] Sets the certificate from file
		--private <FILE>						[Certificate] Sets the private key from file
		--root-string 'ROOT'					[Certificate] Sets the root CA from string input
		--certificate-string 'CERT'       		[Certificate] Sets the certificate from string input
		--private-string 'PRIVATE'				[Certificate] Sets the private key from string input

	generate --type <cred type> --name <cred name> [generate params]
		Generate and set a credential value based on generation parameters.
		-t, --type ['value', 'certificate']		Sets the type of credential to store or generate. (Default: 'value')
		-n, --name 'CRED'						Selects the credential being set

		Generate parameters by [Type]
		-l, --length [4-200]					[Value] Length of generated value (Default: 20)
		--exclude-upper 			        	[Value] Exclude upper alpha characters from generated value
		--exclude-lower 		            	[Value] Exclude lower alpha characters from generated value
		--exclude-number 		            	[Value] Exclude numbers from generated value
		--exclude-special 	  	            	[Value] Exclude special characters from generated value
		--ca 'CA NAME'					     	[Certificate] Name of CA used to sign the generated certificate* (Default: 'default')
		--duration [1-3650]						[Certificate] Valid duration (in days) of the generated certificate* (Default: 365)
		--key-length [2048, 3072, 4096]			[Certificate] Bit length of the generated key (Default: 2048)
		--common-name 'COMMON NAME'				[Certificate] Common name of the generated certificate
		--alternative-name 'ALT NAME'			[Certificate] Alternative name(s) of the generated certificate
		--organization 'ORG'					[Certificate] Organization of the generated certificate
		--organization-unit 'ORG UNIT'			[Certificate] Organization unit of the generated certificate
		--locality 'LOCALITY'					[Certificate] Locality/city of the generated certificate
		--state	'ST'							[Certificate] State/province of the generated certificate
		--country 'CC'							[Certificate] Country of the generated certificate

	delete --name <cred name>
		Delete a credential. 
		-n, --name 'CRED'						Name of credential to delete
```
