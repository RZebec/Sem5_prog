package clientContainer

import (
	"crypto/x509"
	"crypto/tls"
	"io/ioutil"
	"log"
	"strconv"
)

func HttpRequest(baseUrl string, port int, certificatePath string, message string) []byte{
	log.SetFlags(log.Lshortfile)

	cert, err := ioutil.ReadFile(certificatePath)
	if err != nil {
		log.Fatalf("Couldn't load file", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	conf := &tls.Config{
		RootCAs: certPool,
	}

	conn, err := tls.Dial("tcp", baseUrl+":"+strconv.Itoa(port), conf)
	defer conn.Close()

	if err != nil {
		log.Println(err)
	}

	n, err := conn.Write([]byte(message+"\n"))
	if err != nil {
		log.Println(n, err)
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
	}

	println(string(buf[:n]))
	return buf
}
