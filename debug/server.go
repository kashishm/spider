package debug

import (
	"fmt"

	"net"

	"net/http"

	"html/template"
	"log"
)

var infos []processInfo
var t *template.Template

func Start() {
	infos = getPInfos()
	if len(infos) < 1 {
		log.Fatal("No running gauge process found.")
	}
	t, err = template.New("html").Parse(html)
	if err != nil {
		log.Fatalf("Cannot load HTML template. Error: %s", err.Error())
	}

	http.HandleFunc("/", handle)
	http.HandleFunc("/projectRootRequest", projectRootRequest)
	http.HandleFunc("/installationRootRequest", installationRootRequest)
	http.HandleFunc("/libPathRequest", libPathRequest)
	http.HandleFunc("/allStepsRequest", allStepsRequest)
	http.HandleFunc("/allConceptsRequest", allConceptsRequest)
	http.HandleFunc("/specsRequest", specsRequest)
	http.HandleFunc("/stepValueRequest", stepValueRequest)
	http.HandleFunc("/formatSpecsRequest", formatSpecsRequest)
	http.HandleFunc("/refactoringRequest", refactoringRequest)

	p := getFreePort()
	fmt.Printf("Starting server at http://127.0.0.1:%d\n", p)
	log.Fatalf(http.ListenAndServe(fmt.Sprintf(":%d", p), nil).Error())
}

func getFreePort() int {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 0})
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
