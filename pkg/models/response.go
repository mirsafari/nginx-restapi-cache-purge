package models

type PingEndpointResponse struct {
	Status string `json:"status" example:"success"`
}

type CacheEndpointResponseGET struct {
	Status     string  `json:"status" example:"error"`
	DomainName string  `json:"domain" example:"www.domain.tld"`
	Message    string  `json:"msg" example:"cache_folder_not_found"`
	CacheSize  float64 `json:"cacheSize" example:"1024"`
}

type CacheEndpointResponseDELETE struct {
	Status     string `json:"status" example:"error"`
	DomainName string `json:"domain" example:"www.domain.tld"`
	Message    string `json:"msg" example:"cache_folder_not_found"`
}
