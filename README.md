<A name="toc0_1" title="Go clients"/>
# Example clients in Go

##Contents     
**<a href="toc1_1">Encryption</a>**  
**<a href="toc1_2">Http Multi Get</a>**  
**<a href="toc1_3">Convert</a>**  


<A name="toc1_1" title="Encryption" />
## Encryption ##
Simple client used for encrypt/decrypt
<A name="toc1_2" title="Http multi Get" />
## Multi get ##
Kick of goroutines to get multiple http resources using a single client/transport. The main function creates 3 channels (error, done, response) and kicks off a goroutine for each url. Each goroutine creates a new request (all using the same single client/transport) and sends back errors/responses/done on the channels. The main function then selects from the error/done/response/time.After channels repeatedly for the number of goroutines kicked off earlier. ExampleHttpMultiGet retrieves data for 50 stocks from NSE using a single client. 

+ make multiget
+ ./bin/cmd/multiget

<A name="toc1_3" title="Convert" />
## Convert ##
Convert creates conversion tables. ascii.html lists the 256 ascii characters in decimal, hex, unicode and character formats.
