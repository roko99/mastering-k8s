## LAB 01.10

### Entry Level

Using static pods for cluster bootstaping lead to problem that limit of `--max-pods=4` of kubelet application has been reached. To fix this we need to restart kubelet with suitable `--max-pods` value, eg 10.
```bash
sudo PATH=$PATH:/opt/cni/bin:/usr/sbin kubebuilder/bin/kubelet \
    --kubeconfig=/var/lib/kubelet/kubeconfig \
    --config=/var/lib/kubelet/config.yaml \
    --root-dir=/var/lib/kubelet \
    --cert-dir=/var/lib/kubelet/pki \
    --hostname-override=$(hostname) \
    --pod-infra-container-image=registry.k8s.io/pause:3.10 \
    --node-ip=$HOST_IP \
    --cgroup-driver=cgroupfs \
    --max-pods=10  \
    --v=1 &
```
