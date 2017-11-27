package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

func StartWebServer(port string) {
	r := NewRouter()    // NEW
	http.Handle("/", r) // NEW

	logrus.Infof("Starting HTTP service at " + port)
	err := http.ListenAndServe(":"+port, nil) // Goroutine will block here

	if err != nil {
		logrus.Errorln("An error occured starting HTTP listener at port " + port)
		logrus.Errorln("Error: " + err.Error())
	}
}
