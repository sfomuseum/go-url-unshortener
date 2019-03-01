# go-url-unshortener

...

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.11](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

This should still be considered experimental.

## Tools

### unshorten

```
grep expanded_url ./tweet.js | awk '{ print $3 }' | sort | uniq | sed 's/^"//' | sed 's/",$"//' | ./bin/unshorten -stdin -verbose
2019/03/01 11:57:47 Head http://airwaysnews.com/blog/2016/06/09/wow-air-kicks-off-san-francisco-service/%22,: dial tcp: lookup airwaysnews.com: no such host
2019/03/01 11:57:48 http://4sq.com/RfgMa", becomes http://4sq.com/RfgMa%22,
2019/03/01 11:57:48 Head http://CNN.com",: dial tcp: lookup CNN.com",: no such host
2019/03/01 11:57:49 Head http://MahoBeachCam.com",: dial tcp: lookup MahoBeachCam.com",: no such host
2019/03/01 11:57:49 http://500px.com/photo/4751280", becomes https://500px.com:443/photo/4751280%22,
2019/03/01 11:57:49 http://1.usa.gov/1LbxdUe", becomes https://www.transportation.gov/fastlane/women-in-aviation-connect-engage-inspire
2019/03/01 11:57:49 Head http://airwaysnews.com/blog/2015/08/14/american-airlines-to-launch-new-trial-on-uniforms/%22,: dial tcp: lookup airwaysnews.com: no such host
2019/03/01 11:57:49 Head http://airwaysnews.com/blog/2016/04/01/airline-industry-announcements-for-april-1/%22,: dial tcp: lookup airwaysnews.com: no such host
2019/03/01 11:57:49 Head http://airwaysnews.com/blog/2016/04/27/museum-of-flight-completes-final-boeing-247d-flight/%22,: dial tcp: lookup airwaysnews.com: no such host
2019/03/01 11:57:50 http://bit.ly/1rbWi7G", becomes https://www.flysfo.com/museum/aviation-museum-library/collection?field_type_collection_tid_1=1027
2019/03/01 11:57:50 http://bit.ly/1950sConsumer", becomes http://bit.ly/1950sConsumer%22,
2019/03/01 11:57:50 http://bit.ly/16vo1lU", becomes https://www.flysfo.com/museum/exhibitions/souvenirs-tokens-travel07.html
2019/03/01 11:57:51 http://bit.ly/1NvfDZt", becomes https://www.flysfo.com/museum/exhibitions/classic-monsters-kirk-hammett-collection
2019/03/01 11:57:52 http://bit.ly/1TPvJ29", becomes https://www.flysfo.com/museum/about/employment
2019/03/01 11:57:52 http://bit.ly/1RVIYKt", becomes https://www.flysfo.com/museum/public-art-collection?nid=3292
2019/03/01 11:57:52 http://bit.ly/1U7sdhL", becomes https://www.flysfo.com/museum/aviation-museum-library/collection?field_type_collection_tid_1=1025
2019/03/01 11:57:52 http://1.usa.gov/KmedO3", becomes https://www.nga.gov/404status.html
2019/03/01 11:57:52 http://bit.ly/1WLCpRH", becomes https://www.flysfo.com/museum/aviation-museum-library/collection/10319
... and so on
```