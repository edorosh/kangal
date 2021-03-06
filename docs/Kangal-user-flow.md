# Kangal User Flow

We expect users to communicate with Kangal by only using API, which is provided by Kangal Proxy.

>You can import [openapi.json](../openapi.json) file to your Postman and have a collection of requests to Kangal. 

Here is an example of requests users can send to Kangal API to manage their load test.

#### Create 
Create a new load test by making a POST request to Kangal Proxy.
> This example CURL command uses JMeter loadtest type and jmx test file.
> Other load generator types may require other data in request.
<p align="center">  
<img src="jmeter-in-kangal/images/sending_request_postman.png" height="500">
</p>

```
curl -X POST http://${KANGAL_PROXY_ADDRESS}/load-test \
  -H 'Content-Type: multipart/form-data' \
  -F distributedPods=1 \
  -F testFile=@./examples/constant_load.jmx \
  -F testData=@./artifacts/loadtests/testData.csv \
  -F envVars=@./artifacts/loadtests/envVars.csv \
  -F type=JMeter \
  -F overwrite=true
```

#### Check 
Check the status of the load test

```
curl -X GET \
  http://${KANGAL_PROXY_ADDRESS}/load-test/loadtest-name
```

#### Live monitoring
Get logs and monitor your tests. 
Example of monitoring for JMeter described [here](jmeter-in-kangal/Reporting-in-JMeter.md#live-metrics-reporting)
You can also monitor the behavior of your service with your custom tools e.g. Graphite.

```
curl -X GET \
  http://${KANGAL_PROXY_ADDRESS}/load-test/loadtest-name/logs
```

#### Get static report. 
When the test is finished successfully Kangal will save the report, generated by the backend the S3 bucket. 

The report for a particular test can be found by the link `https://${KANGAL_PROXY_ADDRESS}/load-test/loadtest-name/report/`

#### Delete 
Delete your finished load test

```
curl -X DELETE \
  http://${KANGAL_PROXY_ADDRESS}/load-test/loadtest-name
```
