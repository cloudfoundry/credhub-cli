**Pivotal Cloud Foundry CredHub CLI helps you configure and interact with CredHub APIs.**

*Starred commands and parameters are planned, but not yet implemented*

```
Usage: cm [<options>] <command> [<args>]
		--version				    	 		Show version of CLI and API
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
		--ca 'CA NAME'					     	[Certificate] Name of CA used to sign the generated certificate (Default: 'default')
		--duration [1-3650]						[Certificate] Valid duration (in days) of the generated certificate (Default: 365)
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
		
CERTIFICATE AUTHORITY:

	ca-get --name <ca name>
		Get the value and attributes of a CA. 
		-n, --name 'CA'							Name of CA to retrieve

	ca-set --type <ca type> --name <ca name> [set params]
		Set the value and attributes of a CA.
		-t, --type ['root', 'intermediate']		Sets the type of CA to store or generate. (Default: 'root')
		-n, --name 'CA'							Selects the CA being set

		Set parameters by [Type]
		--root	<FILE>							[Intermediate] Sets the root CA from file*
		--certificate <FILE>					[Root/Intermediate] Sets the CA certificate from file
		--private <FILE>						[Root/Intermediate] Sets the CA private key from file
		--root-string 'ROOT'					[Intermediate] Sets the root CA from string input*
		--certificate-string 'CERT'       		[Root/Intermediate] Sets the CA certificate from string input
		--private-string 'PRIVATE'				[Root/Intermediate] Sets the CA private key from string input

	ca-generate --type <ca type> --name <ca name> [generate params]*
		Generate and set a credential value based on generation parameters.
		-t, --type ['value', 'certificate']		Sets the type of credential to store or generate. (Default: 'value')*
		-n, --name 'CRED'						Selects the credential being set*

		Generate parameters by [Type]
		--ca 'CA NAME'					     	[Intermediate] Name of CA used to sign the generated certificate (Default: 'default')*
		--duration [1-3650]						[Root/Intermediate] Valid duration (in days) of the generated certificate (Default: 365)*
		--key-length [2048, 3072, 4096]			[Root/Intermediate] Bit length of the generated key (Default: 2048)*
		--common-name 'COMMON NAME'				[Root/Intermediate] Common name of the generated certificate*
		--alternative-name 'ALT NAME'			[Root/Intermediate] Alternative name(s) of the generated certificate*
		--organization 'ORG'					[Root/Intermediate] Organization of the generated certificate*
		--organization-unit 'ORG UNIT'			[Root/Intermediate] Organization unit of the generated certificate*
		--locality 'LOCALITY'					[Root/Intermediate] Locality/city of the generated certificate*
		--state	'ST'							[Root/Intermediate] State/province of the generated certificate*
		--country 'CC'							[Root/Intermediate] Country of the generated certificate*
```
