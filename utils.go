package utils

import (
	"encoding/json"
	"fmt"
	tls "github.com/Noooste/utls"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"unicode"
)

const (
	Version     = "112"
	UserAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + Version + ".0.0.0 Safari/537.36"
	SecChUa     = "\"Chromium\";v=\"112\", \"Google Chrome\";v=\"112\", \"Not:A-Brand\";v=\"99\""
	BrowserName = "Chrome"
	DeviceType  = "Desktop"
	OsName      = "Windows"
	OsVersion   = "10"
)

func GetSecChUa(version int) string {
	switch version {
	case 106:
		return "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\""
	case 107:
		return "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\""
	case 108:
		return "\"Not?A_Brand\";v=\"8\", \"Chromium\";v=\"108\", \"Google Chrome\";v=\"108\""
	case 109:
		return "\"Not_A Brand\";v=\"99\", \"Google Chrome\";v=\"109\", \"Chromium\";v=\"109\""
	case 110:
		return "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\""
	case 111:
		return "\"Google Chrome\";v=\"111\", \"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"111\""
	case 112:
		return "\"Chromium\";v=\"112\", \"Google Chrome\";v=\"112\", \"Not:A-Brand\";v=\"99\""
	default:
		return SecChUa
	}
}

func GetRandomChromeUserAgent(minVer, maxVer int) (string, string) {
	ver := rand.Intn(maxVer-minVer) + minVer
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + strconv.Itoa(ver) + ".0.0.0 Safari/537.36",
		GetSecChUa(ver)
}

func Must(a any, e error) any {
	if e != nil {
		panic(e)
	}
	return a
}

func SafeGoRoutine(fn func(), recoverFns ...func()) {
	go func() {
		if len(recoverFns) > 0 {
			defer recoverFns[0]()
		}
		defer func() {
			if r := recover(); r != nil {
				log.Print(r)
				debug.PrintStack()
			}
		}()
		fn()
	}()
}

func UrlEncodeMap(query map[string]string) string {
	params := url.Values{}
	for key, value := range query {
		params.Add(key, value)
	}
	return params.Encode()
}

func GetID() string {
	return uuid.New().String()
}

func EscapeQuotes(str string) string {
	return strings.ReplaceAll(str, "\"", "\\\"")
}

func JsonDumps(obj interface{}) string {
	return string(JsonDumpsBytes(obj))
}

func JsonDumpsBytes(obj interface{}) []byte {
	dumped, err := json.Marshal(obj)

	if err != nil {
		return []byte("{}")
	}

	return dumped
}

func UrlEncode(obj any) string {
	r := reflect.ValueOf(obj)
	switch r.Kind() {
	case reflect.Map:
		keys := r.MapKeys()
		var result []string
		for _, key := range keys {
			result = append(result, fmt.Sprintf("%s=%v", key, r.MapIndex(key)))
		}
		return strings.Join(result, "&")

	case reflect.Struct:
		var result []string
		for i := 0; i < r.NumField(); i++ {
			if name, ok := r.Type().Field(i).Tag.Lookup("json"); ok {
				//detect if omitempty is set
				split := strings.Split(name, ",")
				if len(split) > 1 && split[1] == "omitempty" {
					if r.Field(i).IsZero() {
						continue
					}
				}
				result = append(result, fmt.Sprintf("%s=%s", split[0], url.QueryEscape(fmt.Sprintf("%v", r.Field(i)))))
			}

		}
		return strings.Join(result, "&")

	default:
		return ""
	}
}

func FormatProxy(proxy string) string {
	var split = strings.Split(strings.Trim(proxy, "\n\r"), ":")
	if len(split) == 4 {
		return "http://" + split[2] + ":" + split[3] + "@" + split[0] + ":" + split[1]
	} else if len(split) == 2 {
		return "http://" + split[0] + ":" + split[1]
	}
	return proxy
}

func CleanAddress(address string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ := transform.String(t, address)
	var punctuationReplace = regexp.MustCompile(`[[:punct:]]`)
	return punctuationReplace.ReplaceAllString(s, ` `)
}

// GetLastChromeVersion apply the latest Chrome version
// Current Chrome version : 111
func GetLastChromeVersion() *tls.ClientHelloSpec {
	extensions := []tls.TLSExtension{
		&tls.UtlsGREASEExtension{},
		&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
			{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
			{Group: tls.X25519},
		}},
		&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
		&tls.SNIExtension{},
		&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.PKCS1WithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA384,
			tls.PKCS1WithSHA384,
			tls.PSSWithSHA512,
			tls.PKCS1WithSHA512,
		}},
		&tls.UtlsExtendedMasterSecretExtension{},
		&tls.SessionTicketExtension{},
		&tls.SCTExtension{},
		&tls.RenegotiationInfoExtension{},
		&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
			tls.PskModeDHE,
		}},
		&tls.ApplicationSettingsExtension{SupportedProtocols: []string{"h2", "http/1.1"}},
		&tls.CompressCertificateExtension{Algorithms: []tls.CertCompressionAlgo{
			tls.CertCompressionBrotli,
		}},
		&tls.SupportedVersionsExtension{Versions: []uint16{
			tls.GREASE_PLACEHOLDER,
			tls.VersionTLS13,
			tls.VersionTLS12,
		}},
		&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
			tls.GREASE_PLACEHOLDER,
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		}},
		&tls.StatusRequestExtension{},
		&tls.SupportedPointsExtension{SupportedPoints: []byte{
			0x00, // pointFormatUncompressed
		}},
		&tls.UtlsGREASEExtension{},
		&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
	}

	extensionsLength := len(extensions)
	lastTwo := extensionsLength - 2

	rand.Shuffle(extensionsLength, func(i, j int) {
		if i >= lastTwo || j >= lastTwo || i == 0 || j == 0 {
			return
		}
		extensions[i], extensions[j] = extensions[j], extensions[i]
	})

	return &tls.ClientHelloSpec{
		CipherSuites: []uint16{
			tls.GREASE_PLACEHOLDER,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		CompressionMethods: []byte{
			0x00, // compressionNone
		},
		Extensions: extensions,
	}
}
