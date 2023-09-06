package vertexai_test

import (
	"github.com/joho/godotenv"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVertexai(t *testing.T) {
	RegisterFailHandler(Fail)
	godotenv.Load("../.env")

	RunSpecs(t, "Vertexai Suite")
}
