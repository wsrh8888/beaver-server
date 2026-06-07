package coreupdate

import (
	"crypto/sha256"
	"encoding/binary"
	"strings"
)

// RolloutBucket 稳定分桶 0-99，同 device/user 在同一架构下始终相同
func RolloutBucket(appID string, architectureID uint, deviceID, userID string) int {
	key := deviceID
	if userID != "" {
		key = userID
	}
	if key == "" {
		key = "anonymous"
	}
	raw := appID + ":" + itoa(architectureID) + ":" + key
	sum := sha256.Sum256([]byte(raw))
	return int(binary.BigEndian.Uint32(sum[:4]) % 100)
}

func itoa(n uint) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

// CompareVersion 语义化版本比较：latest > current 返回 true
func CompareVersion(latest, current string) bool {
	if current == "" {
		return latest != ""
	}
	return compareParts(splitVersion(latest), splitVersion(current)) > 0
}

// BelowMinVersion current 低于 min 返回 true
func BelowMinVersion(current, min string) bool {
	if min == "" {
		return false
	}
	return compareParts(splitVersion(current), splitVersion(min)) < 0
}

func splitVersion(v string) []int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	v = strings.TrimPrefix(v, "V")
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n := 0
		for _, ch := range p {
			if ch < '0' || ch > '9' {
				break
			}
			n = n*10 + int(ch-'0')
		}
		out = append(out, n)
	}
	return out
}

func compareParts(a, b []int) int {
	max := len(a)
	if len(b) > max {
		max = len(b)
	}
	for i := 0; i < max; i++ {
		ai, bi := 0, 0
		if i < len(a) {
			ai = a[i]
		}
		if i < len(b) {
			bi = b[i]
		}
		if ai > bi {
			return 1
		}
		if ai < bi {
			return -1
		}
	}
	return 0
}

// ResolveInput 发版策略输入
type ResolveInput struct {
	AppID          string
	ArchitectureID uint
	DeviceID       string
	UserID         string
	CurrentVersion string
	StableVersionID uint
	GrayVersionID   uint
	RolloutPercent  uint
	MinVersion      string
	ForceUpdate     bool
	PolicyActive    bool
}

// ResolveResult 目标版本决策
type ResolveResult struct {
	TargetVersionID uint
	InGrayRollout   bool
	ForceUpdate     bool
}

// Resolve 大厂式：正式版 + device/user 比例灰度 + 最低版本强更
func Resolve(in ResolveInput) ResolveResult {
	res := ResolveResult{TargetVersionID: in.StableVersionID}

	if !in.PolicyActive || in.StableVersionID == 0 {
		return res
	}

	inGray := in.GrayVersionID > 0 &&
		in.RolloutPercent > 0 &&
		RolloutBucket(in.AppID, in.ArchitectureID, in.DeviceID, in.UserID) < int(in.RolloutPercent)

	if inGray {
		res.TargetVersionID = in.GrayVersionID
		res.InGrayRollout = true
	}

	force := in.ForceUpdate || BelowMinVersion(in.CurrentVersion, in.MinVersion)
	res.ForceUpdate = force
	return res
}
