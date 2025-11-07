package blob

type CREDENTIAL string

type AzureBlobCredentials struct {
	AccountName   CREDENTIAL
	AccountKey    CREDENTIAL
	ContainerName CREDENTIAL
}

type UploadSASResponse struct {
	SASURL             string `json:"sas_url"`
	PermanentPublicURL string `json:"permanent_public_url"`
	BlobName           string `json:"blob_name"`
}
