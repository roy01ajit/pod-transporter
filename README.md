# pod-transporter

Application to replicate pods running on Kubernetes Namespace A to Kubernetes Namespace B:

Here’s how POST works:
API endpoint expected: ip_address:port/api/v1/pods/replicate
Input data:  {src_namespace: A, target_namespace: B}
Output data: { status: success }
Request method: POST
Posting to the mentioned api end point copy the pods running in src namespace to target namespace
Here’s how GET works:
API endpoint expected: ip_address:port/api/v1/pods/namespace/{namespace}
The above api endpoint should return all the pods running in any namespace
Output data:  {result: [pod1, pod2, pod3, pod4]}
Request method: GET
ip_address:port/api/v1/pods/namespace/B should return pods running in namespace B