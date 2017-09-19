package views

import (
	"net/http"
	"github.com/server-forecaster/model/manager"
	"encoding/json"
	"io"
	"github.com/liip/sheriff"
	"github.com/server-forecaster/model"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	defer model.GetDatabase().Close()
	request.ParseForm()
	identifier := request.Form.Get("identifier")
	password := request.Form.Get("password")
	authManager := manager.CreateAuthenticatorManager()
	if identifier == "" || password == "" {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, `"identifier" and "password" are mandatory fields`)
	} else {
		user := authManager.Login(identifier, password)
		if user == nil {
			writer.WriteHeader(http.StatusUnauthorized)
		} else {
			o := &sheriff.Options{
		Groups: []string{"visible"},
		}
			data, _ := sheriff.Marshal(o, user)
			encodedData, _ := json.Marshal(data)
			writer.WriteHeader(http.StatusAccepted)
			writer.Write(encodedData)
		}
	}
}
