package credhub

import (
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)



func (ch *CredHub) GetAllCertificatesMetadata() ([]credentials.CertificateMetadata, error) {

	return ch.makeGetAllCertificatesRequest()
}


func (ch *CredHub) makeGetAllCertificatesRequest() ([]credentials.CertificateMetadata, error) {
	resp, err := ch.Request(http.MethodGet, "/api/v1/certificates/", nil, nil, true)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)

	dec := json.NewDecoder(resp.Body)
	response := make(map[string][]credentials.CertificateMetadata)


	if err := dec.Decode(&response); err != nil {
		return nil, errors.New("The response body could not be decoded: " + err.Error())
	}

	var ok bool
	var data []credentials.CertificateMetadata


	if data, ok = response["certificates"]; !ok || len(data) == 0 {
		return nil, newCredhubError("response did not contain any credentials", "")
	}

	return data, nil
}

