# TODO
- Name everything appropriately.  hivelocity_* , product_list , etc  Change bare_metal_devices to bare_metal_device

- Squash everything
- Make project public and owned by HV

- Try starting from scratch to work on the project, update all the examples

- Take a quick look at the timeouts for each wait method.
- How are we going to handle situations where `ForceNew` will not work because the configuration that the new server will deploy to is out of stock.
- If someone "force quits their deploy", do we have a way to cleanup the order?
- If someone force exits terraform before it's finished, presumably the server will still be provisioning and eventually move onto their account.  (order and service will be created too)
- How long should terraform wait for a deployment to finish... considering sps fails and servers take a while, we need to find a good balance. (at least in the beginning while we get sps v2 out and solidified)
- If someone updates a ForceNew variable such as OS, which triggers destroying the existing server and redeploying.  How are we going to handle if the stock of that device is gone.  (maybe there is a way to tell terraform with a ForceNew event to keep the same device and we could just run a reload at that point)
- Furthermore, if they change the location, should terraform check that the configuration is in stock at the new location before detroying their existing server.  (maybe there is a way to do validation in terraform)
- Bare metal device datasource

- Update the api-docs.hivelocity.net with Terraform links / instructions and /bare-metal-device endpoints / product/list

- Come up with names for all our products to make life easier for terraform people. Can use the product name instead of the id.

