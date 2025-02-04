package install

import (
	"fmt"
	"github.com/wonderivan/logger"
	"strings"
)

type CleanCluster interface {
	Check
	Clean
}

type JoinNodeAndMaster interface {
	Check
	Send
	Join
}

type Init interface {
	Check
	Send
	PreInit
	Join
	Print
}

type Install interface {
	Check
	Send
	Apply
}

var (
	JoinToken       string
	TokenCaCertHash string
	CertificateKey  string
)

//SealosInstaller is
type SealosInstaller struct {
	Masters []string
	Nodes   []string
	VIP     string
	PkgUrl  string
	Hosts   []string
}

const (
	initMaster0 = "init-master0"
	initMasters = "init-masters"
)

//getCommand("init-master0")
//getCommand("masters")
//getCommand("join")
//get command by versions
func (s *SealosInstaller) getCommand(name string) (cmd string) {
	cmds := make(map[string]string)
	cmds = map[string]string{
		initMaster0: `kubeadm init --config=/root/kubeadm-config.yaml --experimental-upload-certs`,
		initMasters: fmt.Sprintf("kubeadm join %s:6443 --token %s --discovery-token-ca-cert-hash %s --experimental-control-plane --certificate-key %s", s.Masters[0], JoinToken, TokenCaCertHash, CertificateKey),
	}

	if strings.HasPrefix(Version, "v1.15") {
		cmds[initMaster0] = `kubeadm init --config=/root/kubeadm-config.yaml --upload-certs`
	}
	v, ok := cmds[name]
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]fetch command error")
		}
	}()
	if !ok {
		panic(1)
	}
	return v
}

//decode output to join token  hash and key
func decodeOutput(output []byte) {
	s0 := string(output)
	slice := strings.Split(s0, "kubeadm join")
	slice1 := strings.Split(slice[1], "Please note")
	logger.Info("[globals]join command is: %s", slice1[0])
	decodeJoinCmd(slice1[0])
}

//  192.168.0.200:6443 --token 9vr73a.a8uxyaju799qwdjv --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 --experimental-control-plane --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07
func decodeJoinCmd(cmd string) {
	stringSlice := strings.Split(cmd, " ")

	for i, r := range stringSlice {
		switch r {
		case "--token":
			JoinToken = stringSlice[i+1]
		case "--discovery-token-ca-cert-hash":
			TokenCaCertHash = stringSlice[i+1]
		case "--certificate-key":
			CertificateKey = stringSlice[i+1][:64]
		}
	}
}
