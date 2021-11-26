package model

type (
	MessageReceive struct {
		ID   string `json:"id_trigger"`
		Shop string `json:"shop_name"`
		URL  string `json:"url"`
	}

	MessgaeSendDataload struct {
	}

	HealthCheckResponse struct {
		ServiceName string `json:"service_name"`
		Version     string `json:"version"`
		Hostname    string `json:"hostname"`
		Timelife    string `json:"time_life"`
	}
)
