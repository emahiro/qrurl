package jwt

import (
	"os"
	"testing"
)

const (
	testPubKey = `{
		"alg": "RS256",
		"e": "AQAB",
		"ext": true,
		"key_ops": [
		  "verify"
		],
		"kty": "RSA",
		"n": "zunD41HNkfkfoA6EzO2HPZ2Wxk8MPC9uRDfcXClSUewY73V1dYdxw_tSuhVdS-8cLqbVCtxmLoNNnI4oz56Y36i2ouYK1hJS8W2V5XIzfj6_khaSnHrP6iqJnt-Sq2m2PJXJXqyJ-cVCXiD_YCayGaXUG7MMVMQT9GR99UFG8OYmRX-YjgYWFhHITLlF9lD8u02QvfOnh1mVL4_yHOfMQkRyiHiX6KQ3aYR49Fr-b69XVMD57v9Q1nzDlBWpABjadQ7vPvrJEEMGw9v_YnTXicfBze8DkM4iipn92EAEKkn-BaamHolBN2yhAZjSdUYibUmlc3ejZBaMUC_M-eNy6w"
	}`
	testPrivateKey = `{
		"alg": "RS256",
		"d": "BpRgakYbp1yoqMpNZsbW2hq8xKmW1gMcGoc8NqIJwBkl4dd2WLpp37WKN-ieAuAjoUgk1ieUCD6cpTHQEmoUqmgPBrcR-SS0BoFZluY-xPDx3o9hMiClo-ieX7L0UwcViG-q65vI92xSY_PgqwoP98jSKd9TsQ5bJBZd-wKMYvSf65TznL7jkVgrjV2qF1foopLWfyc7bQ1QhkYD4JkSboRbcfNgTbz0BBzk54yEnQIAbt_7P0H6JjdEvahAz1t06MVCspFgbV0aeBDuLOcHXOf40w6lISDM0_a2gMHksgPoXVc8rgZbdEEPq23eL2XvUxz97sAPdlL50SHOjMAcAQ",
		"dp": "QGbsKGjhQzV86wUByMAwXVFZfPN12771lJyAD41BoZQ4OqE-Rj9MwjsrUu-wkahVDs7uk7atqombR0gBRyp9FRpsTloCdYPQwoqE72fHQZbAR-FQtvPlUGxftCNnGCnR5shP46NGWPQawAlUDxNyApMamdow8tEaoryAQl_yeR0",
		"dq": "E62L7SvbU9H6FHimQmdtHpxfRTEm76B3wE0v11n8g5gHaFOZ5oBHsdVjHWdz8aAc84f59-aJVhIIgk3nA-bqAxV7nAX9DchsSbycZeHZQvYt8khwCABo7K906cYw4SCU7Zw995nswHZaABmGNDWwWUTDa_9-eOHhc2G7HskIRwE",
		"e": "AQAB",
		"ext": true,
		"key_ops": [
		  "sign"
		],
		"kty": "RSA",
		"n": "zunD41HNkfkfoA6EzO2HPZ2Wxk8MPC9uRDfcXClSUewY73V1dYdxw_tSuhVdS-8cLqbVCtxmLoNNnI4oz56Y36i2ouYK1hJS8W2V5XIzfj6_khaSnHrP6iqJnt-Sq2m2PJXJXqyJ-cVCXiD_YCayGaXUG7MMVMQT9GR99UFG8OYmRX-YjgYWFhHITLlF9lD8u02QvfOnh1mVL4_yHOfMQkRyiHiX6KQ3aYR49Fr-b69XVMD57v9Q1nzDlBWpABjadQ7vPvrJEEMGw9v_YnTXicfBze8DkM4iipn92EAEKkn-BaamHolBN2yhAZjSdUYibUmlc3ejZBaMUC_M-eNy6w",
		"p": "-XsXM7x6FbXbJYaYMS2QShx9gf8idVBuNrAF_jRAvxzO3V7UYeV5WZwIUvfKg6kQPRJkF3FrLpCRCWaxZI__gnWBNQbSUzLoVth_f9_memqd5RbrJwTgYRx64Mjg2idvOln74iW9D3D4OXEJPIzNhUSBHnzJ_qORDRdRbCKS6ms",
		"q": "1FHqtDX4EBIecuxOWm8UO2fYR3-12reQw8UAG2NgahB6vyLoO8WXu8SHUyM1EPIIE8OKUh-T9gbaKpXuCYWvqwK6A_HiUaiXJTYfSzNMfRF8sJyL8jHFtwSI8eKb6iUtBVwdfwaskMdfxtQuGPJIngNmnhbdCLJ2vR7fZmxcuYE",
		"qi": "-RFyvtfeTIsR8Xfa1AMmG2z080l8Vn8nVzzuqatkt6ewPNLryBlD0iWeGX-kHxMApSu2N8v-dmKKg9lDSpKvxV9jinVMvCTqRkCt73QC6rnO54LPU1LNlA-dTciPTKwtRbTvdfvfJau0T9zRHl6d0qa3F0XT5yRLqf2ngTeERDc"
	}`
)

func TestMain(m *testing.M) {
	os.Setenv("LINE_PUBLIC_KEY_ID", testPubKey)
	os.Setenv("LINE_CHANNEL_ID", "lineChannelId")
	os.Setenv("LINE_PRIVATE_KEY", testPrivateKey)
	m.Run()
	os.Setenv("LINE_PUBLIC_KEY_ID", "")
	os.Setenv("LINE_CHANNEL_ID", "")
	os.Setenv("LINE_PRIVATE_KEY", "")
}

func TestCreateToken(t *testing.T) {
	token, err := CreateToken()
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal("token is empty")
	}
	t.Log("token: ", token)
}
