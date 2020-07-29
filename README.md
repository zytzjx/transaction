### send result to cmc server

### sample
```
{
  "_id": "1c2363a3-9309-4f2e-860f-82146fca60e3",
  "uuid": "1c2363a3-9309-4f2e-860f-82146fca60e3",
  "site": "2",
  "operator": "17543",
  "company": "1",
  "productid": "",
  "sourceModel": "PST_ARD_UNIVERSAL_USB_FD",
  "sourceMake": "Android",
  "errorCode": "1",
  "timeCreated": "2013-05-30T14:37:50.0000000",
  "esnNumber": "99000033137773",
  "portNumber": "1"
  
}
```
<span style="color:red">**Importance:**</span>
* transaction  [HMSET]  
    all info send cmc server
    redis **MUST** include "transaction" 
    transaction **MUST** have some fields as follow:
    * sourceModel
    * sourceMake
    * errorCode
    * esnNumber
    
    Server will reture fail if not include these fields.
