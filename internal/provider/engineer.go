package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
        "errors"
)


// GetEngineers - Returns list of engineers (no auth required)
func (c *Client) GetEngineers() ([]engineerDataSourceModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineers := []engineerDataSourceModel{}
	err = json.Unmarshal(body, &engineers)
	if err != nil {
		return nil, err
	}

	return engineers, nil
}

// GetEngineer - Returns specific engineer (no auth required)
func (c *Client) GetEngineer(engineerId string) (engineerDataSourceModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/id/%s", c.HostURL, engineerId), nil)
	if err != nil {
		return engineerDataSourceModel{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return engineerDataSourceModel{}, err
	}

	engineer := engineerDataSourceModel{}
	err = json.Unmarshal(body, &engineer)
	if err != nil {
		return engineerDataSourceModel{}, err
	}

	return engineer, nil
}



// CreateEngineer - Create new engineer
func (c *Client) CreateEngineer(engineer engineerResourceModel) (*engineerDataSourceModel, error) {
	//Cannot marshal variables that are types.String
	// We have to convert the code into regular go types (string)
	// so that it can be marshalled
	apiEngineer := engineerDataSourceModel{
		Name: engineer.Name.ValueString(),
		Email: engineer.Email.ValueString(),
	}
	rb, err := json.Marshal(apiEngineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newEngineer := engineerDataSourceModel{}
	err = json.Unmarshal(body, &newEngineer)
	if err != nil {
		return nil, err
	}

	return &newEngineer, nil
}

// UpdateEngineer - Updates an engineer
func (c *Client) UpdateEngineer(engineerID string, engineer engineerResourceModel) (*engineerDataSourceModel, error) {
	// Cannot marshal variables that are types.String
	// We have to convert the code into regular go types (string)
	// so that it can be marshalled
	apiEngineer := engineerDataSourceModel{
		Name: engineer.Name.ValueString(),
		Email: engineer.Email.ValueString(),
	}
	rb, err := json.Marshal(apiEngineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.HostURL, engineerID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}


	newEngineer := engineerDataSourceModel{}
	err = json.Unmarshal(body, &newEngineer)
	if err != nil {
		return nil, err
	}

	return &newEngineer, nil
}

// Delete Engineer - Deletes an engineer
func (c *Client) DeleteEngineer(engineerID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/engineers/%s", c.HostURL, engineerID), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

        // Accept body either as raw string or JSON with key "succes"
        if string(body) != "engineer resource deleted" {
                // allow json success message
                var result map[string]string
                if err := json.Unmarshal(body, &result); err == nil {
                        if msg, ok := result["success"]; ok && msg == "engineer resource deleted" {
                                return nil
                        }
                }
                return errors.New(string(body))
        }

	return nil
}







