//go:build windows
// +build windows

package windows_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWindows(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "acceptance/windows")
}
