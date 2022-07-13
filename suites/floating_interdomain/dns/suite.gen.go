// Code generated by gotestmd DO NOT EDIT.
package dns

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
)

type Suite struct {
	base.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/floating_interdomain/dns")
	s.T().Cleanup(func() {
		r.Run(`export KUBECONFIG=$KUBECONFIG1 && kubectl delete service -n kube-system exposed-kube-dns` + "\n" + `export KUBECONFIG=$KUBECONFIG2 && kubectl delete service -n kube-system exposed-kube-dns` + "\n" + `export KUBECONFIG=$KUBECONFIG3 && kubectl delete service -n kube-system exposed-kube-dns`)
	})
	r.Run(`[[ ! -z $KUBECONFIG1 ]]`)
	r.Run(`[[ ! -z $KUBECONFIG2 ]]`)
	r.Run(`[[ ! -z $KUBECONFIG3 ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`kubectl expose service kube-dns -n kube-system --port=53 --target-port=53 --protocol=TCP --name=exposed-kube-dns --type=LoadBalancer`)
	r.Run(`kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
	r.Run(`ip1=$(kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}')` + "\n" + `if [[ $ip1 == *"no value"* ]]; then ` + "\n" + `    ip1=$(kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "hostname"}}')` + "\n" + `    ip1=$(dig +short $ip1 | head -1)` + "\n" + `fi` + "\n" + `echo Selected externalIP: $ip1 for cluster1` + "\n" + `[[ ! -z $ip1 ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`kubectl expose service kube-dns -n kube-system --port=53 --target-port=53 --protocol=TCP --name=exposed-kube-dns --type=LoadBalancer`)
	r.Run(`kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
	r.Run(`ip2=$(kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}')` + "\n" + `if [[ $ip2 == *"no value"* ]]; then ` + "\n" + `    ip2=$(kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "hostname"}}')` + "\n" + `    ip2=$(dig +short $ip2 | head -1)` + "\n" + `fi` + "\n" + `echo Selected externalIP: $ip2 for cluster1` + "\n" + `[[ ! -z $ip2 ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG3`)
	r.Run(`kubectl expose service kube-dns -n kube-system --port=53 --target-port=53 --protocol=TCP --name=exposed-kube-dns --type=LoadBalancer`)
	r.Run(`kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}'`)
	r.Run(`ip3=$(kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "ip"}}')` + "\n" + `if [[ $ip3 == *"no value"* ]]; then ` + "\n" + `    ip3=$(kubectl get services exposed-kube-dns -n kube-system -o go-template='{{index (index (index (index .status "loadBalancer") "ingress") 0) "hostname"}}')` + "\n" + `    ip3=$(dig +short $ip3 | head -1)` + "\n" + `fi` + "\n" + `echo Selected externalIP: $ip3 for cluster1` + "\n" + `[[ ! -z $ip3 ]]`)
	r.Run(`export KUBECONFIG=$KUBECONFIG1`)
	r.Run(`---` + "\n" + `cat > configmap.yaml <<EOF` + "\n" + `apiVersion: v1` + "\n" + `kind: ConfigMap` + "\n" + `metadata:` + "\n" + `  name: coredns` + "\n" + `  namespace: kube-system` + "\n" + `data:` + "\n" + `  Corefile: |` + "\n" + `    .:53 {` + "\n" + `        errors` + "\n" + `        health {` + "\n" + `            lameduck 5s` + "\n" + `        }` + "\n" + `        ready` + "\n" + `        kubernetes cluster.local in-addr.arpa ip6.arpa {` + "\n" + `            pods insecure` + "\n" + `            fallthrough in-addr.arpa ip6.arpa` + "\n" + `            ttl 30` + "\n" + `        }` + "\n" + `        k8s_external my.cluster1` + "\n" + `        prometheus :9153` + "\n" + `        forward . /etc/resolv.conf {` + "\n" + `            max_concurrent 1000` + "\n" + `        }` + "\n" + `        loop` + "\n" + `        reload 5s` + "\n" + `    }` + "\n" + `    my.cluster2:53 {` + "\n" + `      forward . ${ip2}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `    my.cluster3:53 {` + "\n" + `      forward . ${ip3}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `EOF`)
	r.Run(`kubectl apply -f configmap.yaml`)
	r.Run(`cat > custom-configmap.yaml <<EOF` + "\n" + `apiVersion: v1` + "\n" + `kind: ConfigMap` + "\n" + `metadata:` + "\n" + `  name: coredns-custom` + "\n" + `  namespace: kube-system` + "\n" + `data:` + "\n" + `  server.override: |` + "\n" + `    k8s_external my.cluster1` + "\n" + `  proxy2.server: |` + "\n" + `    my.cluster2:53 {` + "\n" + `      forward . ${ip2}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `  proxy3.server: |` + "\n" + `    my.cluster3:53 {` + "\n" + `      forward . ${ip3}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `EOF`)
	r.Run(`kubectl apply -f custom-configmap.yaml`)
	r.Run(`export KUBECONFIG=$KUBECONFIG2`)
	r.Run(`cat > configmap.yaml <<EOF` + "\n" + `apiVersion: v1` + "\n" + `kind: ConfigMap` + "\n" + `metadata:` + "\n" + `  name: coredns` + "\n" + `  namespace: kube-system` + "\n" + `data:` + "\n" + `  Corefile: |` + "\n" + `    .:53 {` + "\n" + `        errors` + "\n" + `        health {` + "\n" + `            lameduck 5s` + "\n" + `        }` + "\n" + `        ready` + "\n" + `        kubernetes cluster.local in-addr.arpa ip6.arpa {` + "\n" + `            pods insecure` + "\n" + `            fallthrough in-addr.arpa ip6.arpa` + "\n" + `            ttl 30` + "\n" + `        }` + "\n" + `        k8s_external my.cluster2` + "\n" + `        prometheus :9153` + "\n" + `        forward . /etc/resolv.conf {` + "\n" + `            max_concurrent 1000` + "\n" + `        }` + "\n" + `        loop` + "\n" + `        reload 5s` + "\n" + `    }` + "\n" + `    my.cluster1:53 {` + "\n" + `      forward . ${ip1}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `    my.cluster3:53 {` + "\n" + `      forward . ${ip3}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `EOF`)
	r.Run(`cat > custom-configmap.yaml <<EOF` + "\n" + `apiVersion: v1` + "\n" + `kind: ConfigMap` + "\n" + `metadata:` + "\n" + `  name: coredns-custom` + "\n" + `  namespace: kube-system` + "\n" + `data:` + "\n" + `  server.override: |` + "\n" + `    k8s_external my.cluster2` + "\n" + `  proxy1.server: |` + "\n" + `    my.cluster1:53 {` + "\n" + `      forward . ${ip1}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `  proxy3.server: |` + "\n" + `    my.cluster3:53 {` + "\n" + `      forward . ${ip3}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `EOF`)
	r.Run(`kubectl apply -f custom-configmap.yaml`)
	r.Run(`kubectl apply -f configmap.yaml`)
	r.Run(`export KUBECONFIG=$KUBECONFIG3`)
	r.Run(`cat > configmap.yaml <<EOF` + "\n" + `apiVersion: v1` + "\n" + `kind: ConfigMap` + "\n" + `metadata:` + "\n" + `  name: coredns` + "\n" + `  namespace: kube-system` + "\n" + `data:` + "\n" + `  Corefile: |` + "\n" + `    .:53 {` + "\n" + `        errors` + "\n" + `        health {` + "\n" + `            lameduck 5s` + "\n" + `        }` + "\n" + `        ready` + "\n" + `        kubernetes cluster.local in-addr.arpa ip6.arpa {` + "\n" + `            pods insecure` + "\n" + `            fallthrough in-addr.arpa ip6.arpa` + "\n" + `            ttl 30` + "\n" + `        }` + "\n" + `        k8s_external my.cluster3` + "\n" + `        prometheus :9153` + "\n" + `        forward . /etc/resolv.conf {` + "\n" + `            max_concurrent 1000` + "\n" + `        }` + "\n" + `        loop` + "\n" + `        reload 5s` + "\n" + `    }` + "\n" + `    my.cluster1:53 {` + "\n" + `      forward . ${ip1}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `    my.cluster2:53 {` + "\n" + `      forward . ${ip2}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `EOF`)
	r.Run(`kubectl apply -f configmap.yaml`)
	r.Run(`cat > custom-configmap.yaml <<EOF` + "\n" + `apiVersion: v1` + "\n" + `kind: ConfigMap` + "\n" + `metadata:` + "\n" + `  name: coredns-custom` + "\n" + `  namespace: kube-system` + "\n" + `data:` + "\n" + `  server.override: |` + "\n" + `    k8s_external my.cluster3` + "\n" + `  proxy1.server: |` + "\n" + `    my.cluster1:53 {` + "\n" + `      forward . ${ip1}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `  proxy2.server: |` + "\n" + `    my.cluster2:53 {` + "\n" + `      forward . ${ip2}:53 {` + "\n" + `        force_tcp` + "\n" + `      }` + "\n" + `    }` + "\n" + `EOF`)
	r.Run(`kubectl apply -f custom-configmap.yaml`)
}
func (s *Suite) Test() {}
