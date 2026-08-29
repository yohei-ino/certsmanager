.. certsmanager usage guide

.. _usage_guide:

=====================================================================
Basic Usage
=====================================================================

The core function of certsmanager is to manage certificates lifecycle.


.. warning::
   All operations performed by this toolkit must be executed in a secure, 
   controlled environment. Never share the CA private key.


Usage Flow
----------

The standard workflow to secure a new server involves calling the 
CA initialization routine first, followed by issuing a certificate for 
the target server.

.. toctree::
   :maxdepth: 2
   :caption: API Details

   api_reference


