package utils

import (
	"fmt"
	"testing"
)

func TestVerifyLicenseString(t *testing.T) {
	licenseString := "eyJkYXRhIjoiZXlKamNtVmhkR1ZrUVhRaU9qRTNOekEwTmpReU16a3NJbVY0Y0dseWVVUmhkR1VpT2lJeU1ESTJMVEF5TFRJd0lpd2laWGh3YVhKNVZHbHRaU0k2TVRjM01UVTBOVFl3TUN3aWNISnZaSFZqZEVsa0lqb3hMQ0p3Y205a2RXTjBUbUZ0WlNJNkltTm9jSEp2WW1VaUxDSnpaWEpwWVd4T2RXMWlaWElpT2lKbU9UYzROVFJtWkMweU1UazBMVGhoWVdJdFlqZzFOaTFtTURReU1EZzBOelkzTUdFaUxDSnpkR0YwZFhNaU9pSmhjSEJ5YjNabFpDSXNJblJwYldWemRHRnRjQ0k2TVRjM01EUTJOREkwTlN3aWRtVnljMmx2YmlJNkl1UzhnZVM0bWlKOSIsInNpZ25hdHVyZSI6ImJaZjVFWkJNZzl3R3BEZ3RPeDlpazJBa3JvVlR5dHk3VGFBLzZoMG5vODMxVHZyQkVJOTRLNFVFZ3FaanI3MGFpQVdQY21EMFNLanI0VTZONFlrWG5Eb3ptRUhyajFHK1hGNFhoYWM3SnBaRUJIanNndEppTXNHc1ZpV1lKcDlOeHlhWXpyT05DRHVkMjRKS212RDM2dHlDNy96YjNNYjJxWmFyOW9aQ1pndEVoYkZKMlQzOVdzWmh4MW5IT3RuL3U0SWt1UUJENEFWUG4zc0dXb3RRNDFCdGlQa3lDMHl1akxib05rM2xIVkpyZnNpY2ZaYnhMakxiYjZubnNQSTloWEQ2QUJ3L3ZCRjVMMDM1SFdsN0w2dmhGemFFdTEyaWdQYVNwUlZRQUFFTmxoYjhRWFlJRUNUaDFnMkcvZDVzbi8vOWlsbWo5RlJTTkhKUWJNcTU5dz09In0="

	// 验证和解密授权字符串
	licenseData, valid, err := VerifyLicenseString(licenseString)
	if err != nil {
		t.Fatalf("VerifyLicenseString failed: %v", err)
	}

	if !valid {
		t.Fatalf("License is invalid")
	}

	// 打印解密结果
	fmt.Println("License verification successful!")
	fmt.Println("Decrypted license data:")
	for key, value := range licenseData {
		fmt.Printf("%s: %v\n", key, value)
	}
}
