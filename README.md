# Purpose
The ARO Wrapper provides the Azure Red Hat Openshift(ARO) RP the ability to install Openshift Clusters with the required configurations for the ARO service.
## Capabilities
The Wrapper uses the Openshift Installer as a library, and inserts custom ignition and manifests as needed. Additionally due to the Library implementation we are able to force bump depdencies in the upstream installer to comply with FedRAMP CVE compliance or apply a critical CVE patch faster than the OpenShift upstream.
## Git Model
The Main branch is not particularly utilized and is earmarked for removal. Every new release branch starting at 4.17 is created from the previous release branch. From there version specific code changes are applied to the release branch. 
We can merge the release branches back into Main, however since Main is not utlized in the workflow there's little need do so and is not reliably done.
Pull requests for dependency bumps need to be done against the most recent release branch and should have justification for the change including a Jira in the ARO project.

