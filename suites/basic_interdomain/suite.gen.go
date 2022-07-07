// Code generated by gotestmd DO NOT EDIT.
package basic_interdomain

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/basic_interdomain/dns"
	"github.com/networkservicemesh/integration-tests/suites/basic_interdomain/loadbalancer"
	"github.com/networkservicemesh/integration-tests/suites/basic_interdomain/nsm"
	"github.com/networkservicemesh/integration-tests/suites/basic_interdomain/spire"
)

type Suite struct {
	base.Suite
	loadbalancerSuite loadbalancer.Suite
	dnsSuite          dns.Suite
	spireSuite        spire.Suite
	nsmSuite          nsm.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite, &s.loadbalancerSuite, &s.dnsSuite, &s.spireSuite, &s.nsmSuite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
}
func (s *Suite) TestNsm_consul() {
	r := s.Runner("../deployments-k8s/examples/nsm_consul")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete deployment counting`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -k nse-auto-scale`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG1 delete -f client/dashboard.yaml`)
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -f networkservice.yaml`)
		r.Run(`consul-k8s uninstall --kubeconfig=$KUBECONFIG2 -auto-approve=true -wipe-data=true`)
	})
	r.Run(`brew tap hashicorp/tap` + "\n" + `brew install hashicorp/tap/consul-k8s`)
	r.Run(`consul-k8s install -config-file=helm-consul-values.yaml -set global.image=hashicorp/consul:1.12.0 -auto-approve --kubeconfig=$KUBECONFIG2`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f networkservice.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f client/dashboard.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f service.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k nse-auto-scale`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f server/counting.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --timeout=5m --for=condition=ready pod -l app=dashboard-nsc`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec -it dashboard -- apk add curl`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec -it dashboard -- curl counting:9001`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 port-forward dashboard 9002:9002`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete deploy counting`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 port-forward dashboard 9002:9002`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f server/counting_nsm.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 port-forward dashboard 9002:9002`)
}
func (s *Suite) TestNsm_istio() {
	r := s.Runner("../deployments-k8s/examples/nsm_istio")
	s.T().Cleanup(func() {
		r.Run(`kubectl --kubeconfig=$KUBECONFIG2 delete -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -k nse-auto-scale ` + "\n" + `kubectl --kubeconfig=$KUBECONFIG1 delete -f productpage/productpage.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete -f networkservice.yaml` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete ns istio-system` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection-` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 delete pods --all`)
	})
	r.Run(`curl -sL https://istio.io/downloadIstioctl | sh -` + "\n" + `export PATH=$PATH:$HOME/.istioctl/bin` + "\n" + `istioctl  install --set profile=minimal -y --kubeconfig=$KUBECONFIG2` + "\n" + `istioctl --kubeconfig=$KUBECONFIG2 proxy-status`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -f networkservice.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 apply -f productpage/productpage.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 apply -k nse-auto-scale`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG2 label namespace default istio-injection=enabled` + "\n" + `` + "\n" + `kubectl --kubeconfig=$KUBECONFIG2 apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/bookinfo/platform/kube/bookinfo.yaml`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 wait --timeout=5m --for=condition=ready pod -l app=deploy/productpage-v1`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/productpage-v1 -c cmd-nsc -- apk add curl`)
	r.Run(`kubectl --kubeconfig=$KUBECONFIG1 exec deploy/productpage-v1 -c cmd-nsc -- curl -s productpage.default:9080/productpage | grep -o "<title>Simple Bookstore App</title>"`)
}
