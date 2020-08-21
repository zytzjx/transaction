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


### CMC define
Transaction:
URL:webserviceserver=http://cmc.futuredial.com/ws/
/insert/

request
method:Post

Format: JSON

sample:

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

If database is OK, the dbid should equal to logid. 

{  
	status: 1/2/3/4,  
	recordid: xxx  
	[error: xxx]  
}   
status: used to indicate the record is result of record  
	1 write to MongoDB success  
	2 write to AMQP server success  
	3 fail  
	4 uuid existing in DB  
error: detail error message comes with status=3  
recordid: equals to uuid reported from client  



### Upload log:   
url: http://cmc-dl.futuredial.com/uploadlog/  


string.Format("-X POST -F \"uuid={0}\" -F \"fileobj=@{1}\" -F \"productid={2}\" -w {3} {4}", uuid, System.IO.Path.GetFileName(logfilename), productid,getCode, _url.ToString()  

curl -X POST -F uuid=Aaaa -F fileobj=@aaa.zip -F productid=5 -w %{http_code} https://httpbin.org/post

success response   
{   
  filename: xxx,   
  done: true/false,  
  uuid: xxxx,   
  md5: xxxx,   
}   

transaction uuid same as upload log uuid
