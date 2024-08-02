<div align="center">

### `macaddr`

MAC Address & Prefix Utility for Go

[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/thatmattlove/go-macaddr/test.yml?style=for-the-badge)](https://github.com/thatmattlove/go-macaddr/actions/workflows/test.yml) [![Docs](https://img.shields.io/badge/godoc-reference-007D9C.svg?style=for-the-badge)](https://pkg.go.dev/github.com/thatmattlove/go-macaddr) [![Code Coverage](https://img.shields.io/codecov/c/github/thatmattlove/go-macaddr?style=for-the-badge)](https://codecov.io/gh/thatmattlove/go-macaddr) [![Go Report Card](https://goreportcard.com/badge/github.com/thatmattlove/go-macaddr?style=for-the-badge)](https://goreportcard.com/report/github.com/thatmattlove/go-macaddr)

</div>

## Installation

```
go get -d github.com/thatmattlove/go-macaddr
```

## Usage

### Single MAC Address

```go
mac, err := macaddr.ParseMACAddress("00:00:5e:00:53:ab")
if err != nil {
    panic(err)
}

mac.ByteString()
// {0,0,94,0,83,171}
mac.Clone()
// MACAddress{0,0,0x5e,0,53,0xab}
mac.Dashes()
// 00-00-5e-00-53-ab
mac.Dots()
// 0000.5e00.53ab
mac.Equal(MACAddress{0,0,0x5e,0,53,0xab})
// true
mac.Format("xxx_xxx_xxx_xxx")
// 000_05e_005_3ab
mac.GEqual(MACAddress{0,0,0x5e,0,53,0xac})
// false
mac.Greater(MACAddress{0,0,0x5e,0,53,0xaa})
// true
mac.Int()
// 1577079723
mac.LEqual(MACAddress{0,0,0x5e,0,53,0xac})
// true
mac.Lesser(MACAddress{0,0,0x5e,0,53,0xaa})
// false
mac.Mask(MACAddress{0xff,0xff,0xff,0xff,0xff,0})
// MACAddress{0,0,0x5e,0,0x53,0}
mac.Next()
// MACAddress{0,0,0x5e,0,0x53,0xac}
mac.NoSeparators()
// 00005e0053ab
mac.OUI()
// 00:00:5e
mac.Previous()
// MACAddress{0,0,0x5e,0,0x53,0xaa}
mac.String()
// 00:00:5e:00:53:ab
```

### MAC Prefix

```go
mac, prefix, err := macaddr.ParseMACPrefix("00:00:5e:00:00:00/24")

if err != nil {
    panic(err)
}

prefix.Count()
// 16777216
prefix.First()
// MACAddress{0,0,0x5e,0,0,0}
iter := prefix.Iter()
for iter.Next() {
    iter.Value()
}
// MACAddress{0,0,0x5e,0,0,0}
// MACAddress{0,0,0x5e,0,0,1}
// MACAddress{0,0,0x5e,0,0,2}
// ...
match, err := prefix.Match("00:00:5e:01:23:45")
match.String()
// 00:00:5e:00:00:00/24
match, err = prefix.Match("00:00:5f:01:23:45")
err.Error()
// '00:00:5f:01:23:45' is not contained within MACPrefix 00:00:5e:00:00:00/24
prefix.OUI()
// 00:00:5e
prefix.PrefixLen()
// 24
prefix.String()
// 00:00:5e:00:00:00/24
prefix.WildcardMask()
// MACPrefix{0,0,0,0xff,0xff,0xff}
```

## Roadmap

Depending on if others find this library useful, EUI-64 support may be added. Please open an issue if you would find this helpful.

![GitHub](https://img.shields.io/github/license/thatmattlove/go-macaddr?color=000&style=for-the-badge)
