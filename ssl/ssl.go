package ssl

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

func ParseCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("证书M解码失败")
	}
	//调用x509的接口
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, err
}

func ParsePrivateKey(data []byte) (any, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("证书M解码失败")
	}
	//调用x509的接口
	PKCS8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		return PKCS8, nil
	}

	PKCS1, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return PKCS1, nil
	}

	ECP, err := x509.ParseECPrivateKey(block.Bytes)
	if err == nil {
		return ECP, nil
	}

	return nil, err
}

func VerifyPrivateKey(cert *x509.Certificate, der []byte) error {
	block, _ := pem.Decode(der)
	if block == nil {
		return errors.New("证书M解码失败")
	}
	//调用x509的接口

	PKCS1, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		if PKCS1.PublicKey.Equal(cert.PublicKey) == false {
			return errors.New("私钥与证书不匹配")
		}
		return nil
	}

	ECP, err := x509.ParseECPrivateKey(block.Bytes)
	if err == nil {
		if ECP.PublicKey.Equal(cert.PublicKey) == false {
			return errors.New("私钥与证书不匹配")
		}
		return nil
	}

	PKCS8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {

		// 根据实际类型使用私钥
		switch key := PKCS8.(type) {
		case *rsa.PrivateKey:
			// 在这里使用 RSA 私钥
			if key.PublicKey.Equal(cert.PublicKey) == false {
				return errors.New("私钥与证书不匹配")
			}
			return nil
		case *ecdsa.PrivateKey:
			// 在这里使用 ECDSA 私钥
			if key.PublicKey.Equal(cert.PublicKey) == false {
				return errors.New("私钥与证书不匹配")
			}
			return nil
		case ed25519.PrivateKey:
			// 在这里使用 Ed25519 私钥
			if key.Equal(cert.PublicKey) == false {
				return errors.New("私钥与证书不匹配")
			}
			return nil
		case ecdh.PrivateKey:
			// 在这里使用 Ed25519 私钥
			if key.Equal(cert.PublicKey) == false {
				return errors.New("私钥与证书不匹配")
			}
			return nil
		default:
			return errors.New("不支持的私钥类型")
		}
	}

	return err
}

func CreateCertificate(commonName string, dnsName []string) ([]byte, []byte, error) {
	maxValue := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, maxValue)

	// 定义：引用IETF的安全领域的公钥基础实施（PKIX）工作组的标准实例化内容
	subject := pkix.Name{
		CommonName: commonName,
	}

	// 设置 SSL证书的属性用途
	certificate509 := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		//IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames: dnsName,
	}

	// 生成指定位数密匙
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	// 生成 SSL公匙
	derBytes, err := x509.CreateCertificate(rand.Reader, &certificate509, &certificate509, &pk.PublicKey, pk)
	if err != nil {
		return nil, nil, err
	}
	certBuf := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// 生成 SSL私匙
	keyBuf := pem.EncodeToMemory(&pem.Block{Type: "RAS PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	return certBuf, keyBuf, err
}
