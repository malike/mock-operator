package controllers

import (
	"bytes"
	"encoding/gob"
	"hash/crc64"
	"strconv"
	"strings"
)

func getServiceName(deploymentName string) string {
	return strings.ToLower(deploymentName) + "-service"
}

func getObjectLabels(deploymentName string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/instance":   deploymentName,
		"app.kubernetes.io/name":       "mock-operator",
		"app.kubernetes.io/managed-by": "malike",
	}
}

func getSelectorLabels(deploymentName string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/instance":   deploymentName,
		"app.kubernetes.io/name":       "demo-operator",
		"app.kubernetes.io/managed-by": "malike",
	}
}
func getHashLabel(labels map[string]string) string {
	return labels["appliedHash"]
}

func convertToByteArray(e any) []byte {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(e)
	if err != nil {
		return nil
	}
	return network.Bytes()
}

func generateHash(s any) string {
	crc64Table := crc64.MakeTable(crc64.ECMA)
	return strconv.FormatUint(crc64.Checksum(convertToByteArray(s), crc64Table), 16)
}
