package ipfilter_test

import (
	"net"
	"testing"

	"github.com/admpub/ip2country"
	"github.com/admpub/ipfilter"
	"github.com/stretchr/testify/assert"
)

func TestSingleIP(t *testing.T) {
	f := ipfilter.New(ipfilter.Options{
		AllowedIPs:     []string{"222.25.118.1"},
		BlockByDefault: true,
		IPDB:           ip2country.Bytes(),
	})
	assert.True(t, f.Allowed("222.25.118.1"), "[1] should be allowed")
	assert.True(t, f.Blocked("222.25.118.2"), "[2] should be blocked")
	assert.True(t, f.NetAllowed(net.IP{222, 25, 118, 1}), "[3] should be allowed")
	assert.True(t, f.NetBlocked(net.IP{222, 25, 118, 2}), "[4] should be blocked")
}

func TestSubnetIP(t *testing.T) {
	f := ipfilter.New(ipfilter.Options{
		AllowedIPs:     []string{"10.0.0.0/16"},
		BlockByDefault: true,
		IPDB:           ip2country.Bytes(),
	})
	assert.True(t, f.Allowed("10.0.0.1"), "[1] should be allowed")
	assert.True(t, f.Allowed("10.0.42.1"), "[2] should be allowed")
	assert.True(t, f.Blocked("10.42.0.1"), "[3] should be blocked")
}

func TestManualCountryCode(t *testing.T) {
	assert.Equal(t, ipfilter.IPToCountry("203.25.111.68"), "AU")
	assert.Equal(t, ipfilter.IPToCountry("216.58.199.67"), "US")
	assert.Equal(t, ipfilter.IPToCountry("116.31.116.51"), "CN")
	assert.Equal(t, ipfilter.IPToCountry("117.175.117.164"), "CN")
}

func TestCountryCodeWhiteList(t *testing.T) {
	f := ipfilter.New(ipfilter.Options{
		AllowedCountries: []string{"AU"},
		BlockByDefault:   true,
		IPDB:             ip2country.Bytes(),
	})
	assert.True(t, f.Allowed("203.25.111.68"), "[1] should be allowed")
	assert.True(t, f.Blocked("216.58.199.67"), "[2] should be blocked")
}

func TestCountryCodeBlackList(t *testing.T) {
	f := ipfilter.New(ipfilter.Options{
		BlockedCountries: []string{"RU", "CN"},
		IPDB:             ip2country.Bytes(),
	})
	assert.True(t, f.Allowed("203.25.111.68"), "[1] AU should be allowed")
	assert.True(t, f.Allowed("216.58.199.67"), "[2] US should be allowed")
	assert.True(t, f.Blocked("116.31.116.51"), "[3] CN should be blocked")
}

func TestDynamicList(t *testing.T) {
	f := ipfilter.New(ipfilter.Options{})
	assert.True(t, f.Allowed("116.31.116.51"), "[1] CN should be allowed")
	f.BlockCountry("CN")
	assert.True(t, f.Blocked("116.31.116.51"), "[1] CN should be blocked")
}
