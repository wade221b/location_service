1. VVIMP : Start your docker desktop
```
cd /project

docker build -t delivery-service:latest .

docker run -p 8000:8000 --name delivery-service -d delivery-service:latest

docker ps

docker logs -f <containerID of delivery-service:latest>
```

//Just in case the above steps dont work (it has been tested heavily on a mac and a windows machine), please follow the following steps

```
cd /project
go mod tidy
cd /src
go run main.go
```
Points to be noted
1. This is a production ready code. Please read the comments everywhere as I have left them there to explain my thought process and future scope of extention for most logical parts.
2. Should set values of "PRICING_SERVICE_HOST" and  "CONSUMER_API" in the environment variables.
3. There is a better way of handling the `distance out of reach` scenario, 
    3.1 Rather than just string comparing, custom error codes can be used, which can be then tested in the handler layer and accordingly the response can be sent.
4. I was writing the test case for the main Business Logic layer `order_price_test.go`, but the mocking of the `dynamic` and `static` api took some time. 
5. instead wrote end to end simulation test cases in the test.py file
    5.1 It needs only the `requests` package installed.
    5.2 Can be run by `python test.py`
6. I did the assignemnt in golang as I have written production grade code in golang and know devops stuff too (or can learn what i dont knwo)
I have written Python production code too and can adapt quickly to the type of project being written in python, and hence i would be more interesed in the python backend role (have written test.py file here), but I would leave it up to you!
