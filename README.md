## blink-crowdstrike
> Use this API specification as a reference for the API endpoints you can use to interact with your Falcon environment. These endpoints support authentication via OAuth2 and interact with detections and network containment. For detailed usage guides and more information about API endpoints that don&#39;t yet support OAuth2, see our [documentation inside the Falcon console](https://falcon.us-2.crowdstrike.com/support/documentation).

## ListDevices
* Search for hosts in your environment by platform, hostname, IP, and other criteria.
<table>
<caption>Action Parameters</caption>
  <thead>
    <tr>
        <th>Param Name</th>
        <th>Param Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
       <tr>
           <td>Filter</td>
           <td>The filter expression that should be used to limit the results.</td>
       </tr>
       <tr>
           <td>Limit</td>
           <td>The maximum records to return. [1-5000]</td>
       </tr>
       <tr>
           <td>Offset</td>
           <td>The offset to start retrieving records from.</td>
       </tr>
       <tr>
           <td>Sort</td>
           <td>The property to sort by (e.g. status.desc or hostname.asc)</td>
       </tr>
    </tr>
  </tbody>
</table>



## SearchAcrossDevices
* Find hosts that have observed a given custom IOC. For details about those hosts, use GET /devices/entities/devices/v1
<table>
<caption>Action Parameters</caption>
  <thead>
    <tr>
        <th>Param Name</th>
        <th>Param Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
       <tr>
           <td>Limit</td>
           <td>The first process to return, where 0 is the latest offset. Use with the offset parameter to manage pagination of results.</td>
       </tr>
       <tr>
           <td>Offset</td>
           <td>The first process to return, where 0 is the latest offset. Use with the limit parameter to manage pagination of results.</td>
       </tr>
       <tr>
           <td>Search By</td>
           <td>
The type of the indicator. Valid types include:

sha256: A hex-encoded sha256 hash string. Length - min: 64, max: 64.

md5: A hex-encoded md5 hash string. Length - min 32, max: 32.

domain: A domain name. Length - min: 1, max: 200.

ipv4: An IPv4 address. Must be a valid IP address.

ipv6: An IPv6 address. Must be a valid IP address.
</td>
       </tr>
       <tr>
           <td>Value</td>
           <td>The string representation of the indicator.</td>
       </tr>
    </tr>
  </tbody>
</table>



## GetDetection
* View information about detections
<table>
<caption>Action Parameters</caption>
  <thead>
    <tr>
        <th>Param Name</th>
        <th>Param Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
       <tr>
           <td>ID</td>
           <td></td>
       </tr>
    </tr>
  </tbody>
</table>



## GetDeviceByID
* Get details on one or more hosts by providing agent IDs (AID). You can get a host&#39;s agent IDs (AIDs) from the /devices/queries/devices/v1 endpoint, the Falcon console or the Streaming API
<table>
<caption>Action Parameters</caption>
  <thead>
    <tr>
        <th>Param Name</th>
        <th>Param Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
       <tr>
           <td>ID</td>
           <td>The host agent IDs used to get details on.</td>
       </tr>
    </tr>
  </tbody>
</table>



## ListDetections
* Search for detection IDs that match a given query
<table>
<caption>Action Parameters</caption>
  <thead>
    <tr>
        <th>Param Name</th>
        <th>Param Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
       <tr>
           <td>Filter</td>
           <td>Filter detections using a query in Falcon Query Language (FQL) An asterisk wildcard * includes all results. Common filter options include: status, device.device_id, max_severity. The full list of valid filter options is extensive. Review it in CrowdStrike&#39;s documentation inside the Falcon console. (https://falcon.crowdstrike.com/documentation/45/falcon-query-language-fql)</td>
       </tr>
       <tr>
           <td>Limit</td>
           <td>The maximum number of detections to return in this response (default: 9999; max: 9999). Use with the `offset` parameter to manage pagination of results.</td>
       </tr>
       <tr>
           <td>Offset</td>
           <td>The first detection to return, where `0` is the latest detection. Use with the `limit` parameter to manage pagination of results.</td>
       </tr>
       <tr>
           <td>Query</td>
           <td>Search all detection metadata for the provided string</td>
       </tr>
       <tr>
           <td>Sort</td>
           <td>Sort detections using these options:

- `first_behavior`: Timestamp of the first behavior associated with this detection
- `last_behavior`: Timestamp of the last behavior associated with this detection
- `max_severity`: Highest severity of the behaviors associated with this detection
- `max_confidence`: Highest confidence of the behaviors associated with this detection
- `adversary_id`: ID of the adversary associated with this detection, if any
- `devices.hostname`: Hostname of the host where this detection was detected

Sort either `asc` (ascending) or `desc` (descending). For example: `last_behavior|asc`</td>
       </tr>
    </tr>
  </tbody>
</table>



## DeleteDevice
* This action will delete a host. After the host is deleted, no new detections for that host will be reported via UI or APIs.
<table>
<caption>Action Parameters</caption>
  <thead>
    <tr>
        <th>Param Name</th>
        <th>Param Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
       <tr>
           <td>Host Agent ID</td>
           <td>The host agent ID (AID) of the host you want to contain. Get an agent ID from a detection, the Falcon console, or the Streaming API.</td>
       </tr>
    </tr>
  </tbody>
</table>

