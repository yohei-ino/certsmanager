.. certsmanager api_reference

.. _api_reference:

============================================================================
API Reference
============================================================================

This document details the core functions available for managing PKI assets.
These functions assume the CA key and certificate have already been 
generated successfully.

Key Structures and Constants
------------------------------

*   **`caKeyFile`**: The filename for the Root CA's private key (`ca.key`).
*   **`caCertFile`**: The filename for the Root CA's public certificate (`ca.crt`).
*   **`serverDir`**: The subdirectory within the project directory where 
    individual server keys and certificates are stored.

.. rubric:: Function Details

The following functions handle the primary lifecycle stages of certificate 
management.

.. function:: InitializeCA(projectName, projectDir)

   Generates the necessary private key and self-signed certificate for the 
   Root CA. This must be run first in any new project directory.

   :param str projectName: A descriptive name for the overall project.
   :param str projectDir: The directory where the CA assets (key/cert) will be stored.
   :raises RuntimeError: If key generation or file writing fails.
   :returns: None on success.
   :rtype: None

.. function:: AddServerCertificate(projectName, serverName, sansList, projectDir)

   Generates a new private key pair for a given server and obtains a signed 
   certificate from the root CA.

   :param str projectName: The project name.
   :param str serverName: The common name (CN) that will appear in the certificate.
   :param str sansList: A comma-separated list of Subject Alternative Names (SANs) (e.g., "api.example.com,localhost"). If empty, no SANs are added.
   :param str projectDir: The directory containing the CA assets.
   :raises RuntimeError: If CA assets cannot be read or if signing fails.
   :returns: None on success.
   :rtype: None


Example Usage Snippet
----------------------

.. code-block:: go

   // 1. Initialize the CA
   err := pki.InitializeCA("MyAwesomeApp", "./keys")
   if err != nil {
       log.Fatalf("Failed to initialize CA: %v", err)
   }

   // 2. Add a new server certificate for api.example.com
   sans := "api.example.com,internal-api"
   err = pki.AddServerCertificate("MyAwesomeApp", "api.example.com", sans, "./keys")
   if err != nil {
       log.Fatalf("Failed to add server cert: %v", err)
   }


