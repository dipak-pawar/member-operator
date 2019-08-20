module github.com/codeready-toolchain/member-operator

require (
	cloud.google.com/go v0.43.0 // indirect
	github.com/codeready-toolchain/api v0.0.0-20190813211537-c2a9a838ce38
	github.com/codeready-toolchain/toolchain-common v0.0.0-20190814120500-ab1e2af882aa
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1
	github.com/go-openapi/spec v0.19.2 // indirect
	github.com/go-openapi/swag v0.19.4 // indirect
	github.com/gobuffalo/envy v1.7.0 // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.4 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/openshift/api v3.9.1-0.20190730142803-0922aa5a655b+incompatible
	github.com/openshift/library-go v0.0.0-20190802153730-cab1b771d4dc
	github.com/operator-framework/operator-sdk v0.10.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.3 // indirect
	github.com/redhat-cop/operator-utils v0.0.0-20190530184149-66ee667a40b2
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa // indirect
	golang.org/x/tools v0.0.0-20190807223507-b346f7fd45de // indirect
	google.golang.org/grpc v1.22.0 // indirect
	k8s.io/api v0.0.0-20190718062839-c8a0b81cb10e
	k8s.io/apiextensions-apiserver v0.0.0-20190718063925-2249b0201a0a // indirect
	k8s.io/apimachinery v0.0.0-20190717022731-0bb8574e0887
	k8s.io/apiserver v0.0.0-20190111033246-d50e9ac5404f // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/code-generator v0.0.0-20190717022600-77f3a1fe56bb
	k8s.io/gengo v0.0.0-20190327210449-e17681d19d3a
	k8s.io/klog v0.3.3
	k8s.io/kube-openapi v0.0.0-20190718094010-3cf2ea392886
	sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools v0.1.12
	sigs.k8s.io/kubefed v0.1.0-rc2
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20181117043124-c2090bec4d9b
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20180711000925-0cf8f7e6ed1d
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)
