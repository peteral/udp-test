#!/bin/bash

helm install cilium cilium/cilium --version 1.14.1 \
   --namespace kube-system \
   --set operator.replicas=1 \
   --set egressGateway.enabled=true \
   --set egressGateway.installRoutes=true \
   --set bpf.masquerade=true \
   --set kubeProxyReplacement=strict \
   --set l7Proxy=false \
   --set debug.enabled=true

kubectl label nodes nb-lpetera-u egress-node=true