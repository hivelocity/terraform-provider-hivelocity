# TODO
- Try starting from scratch to work on the project, update all the examples

# Questions that need answering:

1. How are we going to handle situations where `ForceNew` will not work because the configuration that the new server will deploy to is out of stock.
2. If someone "force quits their deploy", do we have a way to cleanup the order?
3. If someone force exits terraform before it's finished, presumably the server will still be provisioning and eventually move onto their account.  (order and service will be created too). How does everyone else handle this.
4. If someone updates a ForceNew variable such as OS, which triggers destroying the existing server and redeploying.  How are we going to handle if the stock of that device is gone.  (maybe there is a way to tell terraform with a ForceNew event to keep the same device and we could just run a reload at that point)
5. Furthermore, if they change the location, should terraform check that the configuration is in stock at the new location before detroying their existing server.  (maybe there is a way to do validation in terraform)
6. Do we need a way to "lock" sensitive fields on devices (and other resources) that can be changed from the portal?

# TODO For other departments
1. Come up with names for all our products to make life easier for terraform people. Can use the product name instead of the id.
