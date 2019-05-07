// This file is generated - do not edit!

package platformsh

import "fmt"

type AccessLevel uint8

const (
	AccessLevelViewer AccessLevel = iota
	AccessLevelContributor
	AccessLevelAdmin
	totalAccessLevels
)

var accessLevels = [totalAccessLevels]string{
	"viewer",
	"contributor",
	"admin",
}

var accessLevelsMap = map[string]AccessLevel{
	"viewer":      AccessLevelViewer,
	"contributor": AccessLevelContributor,
	"admin":       AccessLevelAdmin,
}

func NewAccessLevel(name string) (AccessLevel, error) {
	if v, ok := accessLevelsMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown AccessLevel name %q", name)
}

func (v AccessLevel) String() string {
	if v < totalAccessLevels {
		return accessLevels[v]
	}

	return fmt.Sprintf("unknown AccessLevel value %02x", uint8(v))
}

func (v *AccessLevel) UnmarshalText(text []byte) (err error) {
	*v, err = NewAccessLevel(string(text))
	return err
}

func (v AccessLevel) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

type AccessType uint8

const (
	AccessTypeSSH AccessType = iota
	totalAccessTypes
)

var accessTypes = [totalAccessTypes]string{
	"ssh",
}

var accessTypesMap = map[string]AccessType{
	"ssh": AccessTypeSSH,
}

func NewAccessType(name string) (AccessType, error) {
	if v, ok := accessTypesMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown AccessType name %q", name)
}

func (v AccessType) String() string {
	if v < totalAccessTypes {
		return accessTypes[v]
	}

	return fmt.Sprintf("unknown AccessType value %02x", uint8(v))
}

func (v *AccessType) UnmarshalText(text []byte) (err error) {
	*v, err = NewAccessType(string(text))
	return err
}

func (v AccessType) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

type ApplicationMount uint8

const (
	ApplicationMountLocal ApplicationMount = iota
	ApplicationMountTemp
	ApplicationMountService
	totalApplicationMounts
)

var applicationMounts = [totalApplicationMounts]string{
	"local",
	"tmp",
	"service",
}

var applicationMountsMap = map[string]ApplicationMount{
	"local":   ApplicationMountLocal,
	"tmp":     ApplicationMountTemp,
	"service": ApplicationMountService,
}

func NewApplicationMount(name string) (ApplicationMount, error) {
	if v, ok := applicationMountsMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown ApplicationMount name %q", name)
}

func (v ApplicationMount) String() string {
	if v < totalApplicationMounts {
		return applicationMounts[v]
	}

	return fmt.Sprintf("unknown ApplicationMount value %02x", uint8(v))
}

func (v *ApplicationMount) UnmarshalText(text []byte) (err error) {
	*v, err = NewApplicationMount(string(text))
	return err
}

func (v ApplicationMount) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

type ServiceSize uint8

const (
	ServiceSizeAuto ServiceSize = iota
	ServiceSizeSmall
	ServiceSizeMedium
	ServiceSizeLarge
	ServiceSizeExtraLarge
	ServiceSizeDoubleExtraLarge
	ServiceSizeQuadrupleExtraLarge
	totalServiceSizes
)

var serviceSizes = [totalServiceSizes]string{
	"AUTO",
	"S",
	"M",
	"L",
	"XL",
	"2XL",
	"4XL",
}

var serviceSizesMap = map[string]ServiceSize{
	"AUTO": ServiceSizeAuto,
	"S":    ServiceSizeSmall,
	"M":    ServiceSizeMedium,
	"L":    ServiceSizeLarge,
	"XL":   ServiceSizeExtraLarge,
	"2XL":  ServiceSizeDoubleExtraLarge,
	"4XL":  ServiceSizeQuadrupleExtraLarge,
}

func NewServiceSize(name string) (ServiceSize, error) {
	if v, ok := serviceSizesMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown ServiceSize name %q", name)
}

func (v ServiceSize) String() string {
	if v < totalServiceSizes {
		return serviceSizes[v]
	}

	return fmt.Sprintf("unknown ServiceSize value %02x", uint8(v))
}

func (v *ServiceSize) UnmarshalText(text []byte) (err error) {
	*v, err = NewServiceSize(string(text))
	return err
}

func (v ServiceSize) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

type SocketFamily uint8

const (
	SocketFamilyTCP SocketFamily = iota
	SocketFamilyUnix
	totalSocketFamilies
)

var socketFamilies = [totalSocketFamilies]string{
	"tcp",
	"unix",
}

var socketFamiliesMap = map[string]SocketFamily{
	"tcp":  SocketFamilyTCP,
	"unix": SocketFamilyUnix,
}

func NewSocketFamily(name string) (SocketFamily, error) {
	if v, ok := socketFamiliesMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown SocketFamily name %q", name)
}

func (v SocketFamily) String() string {
	if v < totalSocketFamilies {
		return socketFamilies[v]
	}

	return fmt.Sprintf("unknown SocketFamily value %02x", uint8(v))
}

func (v *SocketFamily) UnmarshalText(text []byte) (err error) {
	*v, err = NewSocketFamily(string(text))
	return err
}

func (v SocketFamily) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

type SocketProtocol uint8

const (
	SocketProtocolHTTP SocketProtocol = iota
	SocketProtocolFastCGI
	SocketProtocolUWSGI
	totalSocketProtocols
)

var socketProtocols = [totalSocketProtocols]string{
	"http",
	"fastcgi",
	"uwsgi",
}

var socketProtocolsMap = map[string]SocketProtocol{
	"http":    SocketProtocolHTTP,
	"fastcgi": SocketProtocolFastCGI,
	"uwsgi":   SocketProtocolUWSGI,
}

func NewSocketProtocol(name string) (SocketProtocol, error) {
	if v, ok := socketProtocolsMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown SocketProtocol name %q", name)
}

func (v SocketProtocol) String() string {
	if v < totalSocketProtocols {
		return socketProtocols[v]
	}

	return fmt.Sprintf("unknown SocketProtocol value %02x", uint8(v))
}

func (v *SocketProtocol) UnmarshalText(text []byte) (err error) {
	*v, err = NewSocketProtocol(string(text))
	return err
}

func (v SocketProtocol) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}
