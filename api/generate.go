package api

import (
	"github.com/cloudfoundry-incubator/credhub-cli/actions"
	"github.com/cloudfoundry-incubator/credhub-cli/client"
	"github.com/cloudfoundry-incubator/credhub-cli/models"
	"github.com/cloudfoundry-incubator/credhub-cli/repositories"
)

func (a *Api) Generate(
	credentialIdentifier string, credentialType string, noOverwrite bool,
	username string, length int, excludeUpper bool, excludeLower bool, excludeNumber bool, includeSpecial bool,
	commonName string, alternativeName []string,
	organization string, organizationUnit string, locality string, state string, country string,
	keyLength int, keyUsage []string, extendedKeyUsage []string, sshComment string,
	duration int, ca string, isCA bool, selfSign bool,
) (models.CredentialResponse, error) {

	repository := repositories.NewCredentialRepository(client.NewHttpClient(*a.Config))

	parameters := models.GenerationParameters{
		IncludeSpecial:   includeSpecial,
		ExcludeNumber:    excludeNumber,
		ExcludeUpper:     excludeUpper,
		ExcludeLower:     excludeLower,
		Length:           length,
		CommonName:       commonName,
		Organization:     organization,
		OrganizationUnit: organizationUnit,
		Locality:         locality,
		State:            state,
		Country:          country,
		AlternativeName:  alternativeName,
		ExtendedKeyUsage: extendedKeyUsage,
		KeyUsage:         keyUsage,
		KeyLength:        keyLength,
		Duration:         duration,
		Ca:               ca,
		SelfSign:         selfSign,
		IsCA:             isCA,
		SshComment:       sshComment,
	}

	var value *models.ProvidedValue
	if len(username) > 0 {
		value = &models.ProvidedValue{
			Username: username,
		}
	}

	action := actions.NewAction(repository, a.Config)
	request := client.NewGenerateCredentialRequest(*a.Config, credentialIdentifier, parameters, value, credentialType, !noOverwrite)
	credential, err := action.DoAction(request, credentialIdentifier)
	if err != nil {
		return models.CredentialResponse{}, err
	}

	return credential.(models.CredentialResponse), err
}

func (a *Api) GeneratePassword(
	credentialIdentifier string, noOverwrite bool,
	length int, excludeUpper bool, excludeLower bool, excludeNumber bool, includeSpecial bool,
) (models.CredentialResponse, error) {
	return a.Generate(
		credentialIdentifier, "password", noOverwrite,
		"", length, excludeUpper, excludeLower, excludeNumber, includeSpecial,
		"", []string{},
		"", "", "", "", "",
		0, []string{}, []string{}, "",
		0, "", false, false,
	)
}

func (a *Api) GenerateUser(
	credentialIdentifier string, noOverwrite bool,
	username string, length int, excludeUpper bool, excludeLower bool, excludeNumber bool, includeSpecial bool,
) (models.CredentialResponse, error) {
	return a.Generate(
		credentialIdentifier, "user", noOverwrite,
		username, length, excludeUpper, excludeLower, excludeNumber, includeSpecial,
		"", []string{},
		"", "", "", "", "",
		0, []string{}, []string{}, "",
		0, "", false, false,
	)
}

func (a *Api) GenerateCertificate(
	credentialIdentifier string, noOverwrite bool,
	commonName string, alternativeName []string,
	organization string, organizationUnit string, locality string, state string, country string,
	keyLength int, keyUsage []string, extendedKeyUsage []string,
	duration int, ca string, isCA bool, selfSign bool,
) (models.CredentialResponse, error) {
	return a.Generate(
		credentialIdentifier, "user", noOverwrite,
		"", 0, false, false, false, false,
		commonName, alternativeName,
		organization, organizationUnit, locality, state, country,
		keyLength, keyUsage, extendedKeyUsage, "",
		duration, ca, isCA, selfSign,
	)
}

func (a *Api) GenerateRsa(
	credentialIdentifier string, noOverwrite bool,
	keyLength int,
) (models.CredentialResponse, error) {
	return a.Generate(
		credentialIdentifier, "rsa", noOverwrite,
		"", 0, false, false, false, false,
		"", []string{},
		"", "", "", "", "",
		keyLength, []string{}, []string{}, "",
		0, "", false, false,
	)
}

func (a *Api) GenerateSsh(
	credentialIdentifier string, noOverwrite bool,
	keyLength int, sshComment string,
) (models.CredentialResponse, error) {
	return a.Generate(
		credentialIdentifier, "ssh", noOverwrite,
		"", 0, false, false, false, false,
		"", []string{},
		"", "", "", "", "",
		keyLength, []string{}, []string{}, sshComment,
		0, "", false, false,
	)
}
