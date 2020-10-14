# TOOO
1. add lifecycle to resource.  And make it not "destroy" if the forceNew is running and the new product_id does not exist in stock.

# Questions that need answering:

1. How are we going to handle situations where `ForceNew` will not work because the configuration that the new server will deploy to is out of stock.
1.5. Furthermore, if they change the location, should terraform check that the configuration is in stock at the new location before detroying their existing server.  (maybe there is a way to do validation in terraform)

See #1 above.

2. If someone updates a ForceNew variable such as OS, which triggers destroying the existing server and redeploying.  How are we going to handle if the stock of that device is gone.  (maybe there is a way to tell terraform with a ForceNew event to keep the same device and we could just run a reload at that point)

Address this when sps v2 is in production.

# TODO For other departments
1. Come up with names for all our products to make life easier for terraform people. Can use the product name instead of the id.
