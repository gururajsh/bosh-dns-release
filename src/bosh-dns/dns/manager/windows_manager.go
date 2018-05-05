package manager

import (
	"fmt"
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

const prependDNSServer = `
param ($DNSAddress = $(throw "DNSAddress parameter is required."))

$ErrorActionPreference = "Stop"

function DnsServers($interface) {
  return (Get-DnsClientServerAddress -InterfaceAlias $interface -AddressFamily ipv4 -ErrorAction Stop).ServerAddresses
}

try {
  # identify our interface
  [array]$routeable_interfaces = Get-WmiObject Win32_NetworkAdapterConfiguration | Where { $_.IpAddress -AND ($_.IpAddress | Where { $addr = [Net.IPAddress] $_; $addr.AddressFamily -eq "InterNetwork" -AND ($addr.address -BAND ([Net.IPAddress] "255.255.0.0").address) -ne ([Net.IPAddress] "169.254.0.0").address }) }
  $interface = (Get-WmiObject Win32_NetworkAdapter | Where { $_.DeviceID -eq $routeable_interfaces[0].Index }).netconnectionid

  # avoid prepending if we happen to already be at the top to try and avoid races
  [array]$servers = DnsServers($interface)
  if($servers[0] -eq $DNSAddress) {
    Exit 0
  }

  Set-DnsClientServerAddress -InterfaceAlias $interface -ServerAddresses (,$DNSAddress + $servers)

  # read back the servers in case set silently failed
  [array]$servers = DnsServers($interface)
  if($servers[0] -ne $DNSAddress) {
      Write-Error "Failed to set '${DNSAddress}' as the first dns client server address"
  }
} catch {
  $Host.UI.WriteErrorLine($_.Exception.Message)
  Exit 1
}

Exit 0
`

type windowsManager struct {
	address        string
	runner         boshsys.CmdRunner
	fs             boshsys.FileSystem
	adapterFetcher AdapterFetcher
}

//go:generate counterfeiter . AdapterFetcher

type AdapterFetcher interface {
	Adapters() ([]Adapter, error)
}

func NewWindowsManager(address string, runner boshsys.CmdRunner, fs boshsys.FileSystem, adapterFetcher AdapterFetcher) *windowsManager {
	return &windowsManager{address: address, runner: runner, fs: fs, adapterFetcher: adapterFetcher}
}

func (manager *windowsManager) SetPrimary() error {
	servers, err := manager.Read()
	if err != nil {
		return err
	}

	if len(servers) > 0 && servers[0] == manager.address {
		return nil
	}

	scriptName, err := manager.writeScript("prepend-dns-server", prependDNSServer)
	if err != nil {
		return bosherr.WrapError(err, "Creating prepend-dns-server.ps1")
	}
	defer manager.fs.RemoveAll(filepath.Dir(scriptName))

	_, _, _, err = manager.runner.RunCommand("powershell.exe", scriptName, manager.address)
	if err != nil {
		return bosherr.WrapError(err, "Executing prepend-dns-server.ps1")
	}

	return nil
}

func (manager *windowsManager) Read() ([]string, error) {
	adapter, err := manager.getPrimaryAdapter()
	if err != nil {
		return nil, bosherr.WrapError(err, "Getting list of current DNS Servers")
	}

	return adapter.DNSServerAddresses, nil
}

func (manager *windowsManager) writeScript(name, contents string) (string, error) {
	dir, err := manager.fs.TempDir(name)
	if err != nil {
		return "", err
	}

	scriptName := filepath.Join(dir, name+".ps1")
	err = manager.fs.WriteFileString(scriptName, contents)
	if err != nil {
		return "", err
	}

	err = manager.fs.Chmod(scriptName, 0700)
	if err != nil {
		return "", err
	}

	return scriptName, nil
}

const (
	IfOperStatusUp         uint32 = 1
	IfTypeSoftwareLoopback uint32 = 24
	IfTypeTunnel           uint32 = 131
)

func (manager *windowsManager) getPrimaryAdapter() (Adapter, error) {
	adapters, err := manager.adapterFetcher.Adapters()
	if err != nil {
		return Adapter{}, err
	}

	var candidateAdapters []Adapter

	for _, adapter := range adapters {
		if adapter.IfType == IfTypeSoftwareLoopback || adapter.IfType == IfTypeTunnel {
			continue
		} else if adapter.OperStatus != IfOperStatusUp {
			continue
		}

		candidateAdapters = append(candidateAdapters, adapter)
	}

	if len(candidateAdapters) == 1 {
		return candidateAdapters[0], nil
	}

	for _, adapter := range candidateAdapters {
		for _, unicastAddress := range adapter.UnicastAddresses {
			if unicastAddress == manager.address {
				return adapter, nil
			}
		}
	}

	return Adapter{}, fmt.Errorf("Unable to find primary adapter for %s", manager.address)
}

type Adapter struct {
	IfType             uint32
	OperStatus         uint32
	UnicastAddresses   []string
	DNSServerAddresses []string
}
