package Logger

import (
	"crypto/rsa"
	"CTng/crypto"
	"net/http"
	"crypto/x509"
	"CTng/config"
	"CTng/CA"
	//"fmt"
)

type Logger_public_config struct {
	All_CA_URLs     []string
	All_Logger_URLs []string
	MMD             int
	MRD             int
	Http_vers       []string
}

type Logger_private_config struct {
	Signer                 string
	Port                   string
	CAlist 		           []string
	Monitorlist 		   []string
	Gossiperlist 		   []string
}


type LoggerContext struct {
	Client *http.Client
	SerialNumber int
	Logger_public_config *Logger_public_config
	Logger_private_config *Logger_private_config
	Logger_crypto_config *crypto.CryptoConfig
	PublicKey rsa.PublicKey
	PrivateKey rsa.PrivateKey
	CurrentPrecertPool *CA.CertPool
	PrecertStorage *PrecertStorage
	OnlinePeriod int
}


type PrecertStorage struct {
	PrecertPools map[string] *CA.CertPool
}

//check if an item is in a list
func inList(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func Verifyprecert (precert x509.Certificate, ctx LoggerContext) bool {
	issuer := precert.Issuer.CommonName
	//check if issuer is in CAlist
	if !inList(issuer, ctx.Logger_private_config.CAlist) {
		return false
	}
	//retrieve the public key of the issuer
	issuerPublicKey := ctx.Logger_crypto_config.SignaturePublicMap[crypto.CTngID(issuer)]
	//retrieve the signature of the precert
	signature := precert.Signature
	rsasig := new(crypto.RSASig)
	(*rsasig).Sig = signature
	(*rsasig).ID = crypto.CTngID(issuer)
	//check if the signature is valid
	if err := crypto.RSAVerify(precert.RawTBSCertificate, *rsasig, &issuerPublicKey); err != nil {
		return false
	}	
	return true
}



// initialize Logger context
func InitializeLoggerContext(public_config_path string,private_config_file_path string,crypto_config_path string) *LoggerContext{
	// Load public config from file
	pubconf := new(Logger_public_config)
	config.LoadConfiguration(&pubconf, public_config_path)
	// Load private config from file
	privconf := new(Logger_private_config)
	config.LoadConfiguration(&privconf, private_config_file_path)
	// Load crypto config from file
	cryptoconfig, err := crypto.ReadCryptoConfig(crypto_config_path)
	if err != nil {
		//fmt.Println("read crypto config failed")
	}
	// Initialize Logger Context
	loggerContext := &LoggerContext{
		SerialNumber: 0,
		Logger_public_config: pubconf,
		Logger_private_config: privconf,
		Logger_crypto_config: cryptoconfig,
		PublicKey:  cryptoconfig.SignaturePublicMap[cryptoconfig.SelfID],
		PrivateKey: cryptoconfig.RSAPrivateKey,
		CurrentPrecertPool: CA.NewCertPool(),
		PrecertStorage: &PrecertStorage{PrecertPools: make(map[string] *CA.CertPool)},
		OnlinePeriod: 0,
	}
	// Initialize http client
	tr := &http.Transport{}
	loggerContext.Client = &http.Client{
		Transport: tr,
	}
	return loggerContext
}

func GenerateLogger_private_config_template() *Logger_private_config{
	return &Logger_private_config{
		Signer: "",
		Port: "",
		CAlist: []string{},
		Monitorlist: []string{},
		Gossiperlist: []string{},
	}
}

func GenerateLogger_public_config_template() *Logger_public_config{
	return &Logger_public_config{
		All_CA_URLs: []string{},
		All_Logger_URLs: []string{},
		MMD: 0,
		MRD: 0,
		Http_vers: []string{},
	}
}

func GenerateLogger_crypto_config_template() *crypto.StoredCryptoConfig{
	return &crypto.StoredCryptoConfig{
		SelfID: crypto.CTngID("0"),
		Threshold: 0,
		N: 0,
		HashScheme: 0,
		SignScheme: "",
		ThresholdScheme: "",
		SignaturePublicMap: crypto.RSAPublicMap{},
		RSAPrivateKey: rsa.PrivateKey{},
		ThresholdPublicMap: map[string][]byte{},
		ThresholdSecretKey: []byte{},
	}
}