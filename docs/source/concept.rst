.. certsmanager concept

.. _concept:

============================================================================
Concepts
============================================================================


The toolkit operates on the principle of issuing digitally signed 
certificates to servers. This process involves three main steps:

1.  **CA Initialization**: Creating the root Certificate Authority (CA).
    This step generates the private key and the self-signed certificate 
    that will be used to vouch for the identity of all other services.
2.  **Server Key Generation**: Generating a unique private key pair 
    for each service (e.g., an API gateway, a backend server).
3.  **Certificate Signing**: Using the CA's private key to sign the 
    server's public key, thereby issuing a trusted certificate.


*Note: Actual diagrams or usage screenshots would be inserted here.*


.. rubric:: Certificate Lifecycle

The certificate issuance process in certsmanager is outlined below.

.. mermaid::

    sequenceDiagram
        autonumber
        actor User as User
        participant CA as Root CA (InitializeCA)
        participant Server as Server (AddServerCertificate)

        User->>CA: InitializeCA(projectName, projectDir)
        Note over CA: Generate private key and self-signed cert for Root CA
        
        User->>Server: AddServerCertificate(..., sansList)
        Server->>CA: Request signature
        CA-->>Server: Issue signed certificate


