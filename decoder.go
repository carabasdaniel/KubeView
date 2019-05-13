package main

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"k8s.io/apimachinery/pkg/util/yaml"

	"k8s.io/client-go/kubernetes/scheme"
)

func Decode(fileContent []byte) []interface{} {
	b := bufio.NewReader(bytes.NewReader(fileContent))
	r := yaml.NewYAMLReader(b)
	var decodedResult []interface{}
	for {
		doc, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		d := scheme.Codecs.UniversalDeserializer()
		obj, _, err := d.Decode(doc, nil, nil)
		if err != nil {
			log.Fatalf("could not decode yaml: %s\n%s", fileContent, err)
		}
		decodedResult = append(decodedResult, obj)
	}
	return decodedResult
}
