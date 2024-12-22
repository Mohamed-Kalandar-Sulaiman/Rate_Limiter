# Ideas
## Supported Algorithms 
    1. Sliding window
    2. Fixed window


- ### GRPC Request
    1. service_name ~ str - Identifier of the service
    2. action_name ~ str - Identifier of the action
    3. config_name ~ str - Identifier of the config plan
    3. uid ~ str - actOR 
    4. oid ~ str - actON
 


- ### GRPC Response
    1. is_rate_limited ~ bool 
    2. limit ~ int
    3. remaining ~ int
    4. reset_after ~ int
    5. reset ~ int


- # RateLimitFunction
    ## Input Args
        - key (service_id:action_id:config_id:algorithm[:2]:uid:oid)
        - algorithm
        - unit
        - request_per_unit
    ## Output
        - is_rate_limited
        - limit
        - remaining
        - reset_after
        - reset

    ## Algorithm
        1. Generate key from the input args
        2. Fetch config data like (algorithm , unit , request_per_unit) from in-memory config object,
         using key service_name:action_name:config_name
        3. Using this unique key we can get the following data
            1. service_id
            2. action_id
            3. config_id
            4. algorithm
            5.  

## Flows

- Each service makes an grpc call to rate limit the user or ip for an action
- The action can either be configured as [ absolute / relative ]
- The action can be configured either as sliding [window / fixed] window for accuracy





## GRPC Endpoints
1. RateLimit
    - serviceId
    - 

## REST API 

1. Service [Microservice]
    - name ~ Human friendly unique identifier provided by each micro service 
    - id ~ auto incremented system generated short id to be used in redis Key 

API ENDPOINTS :-
    Create Service
    Update Service
    Get Service
    Delete Service


2. Actions 
    - service_id ~ foreignKey
    - service_name ~ stored as back ref for human friendly
    - name ~ unique name of the action provided by service team
    - id ~ auto incremented system generated short id to be used in redis Key 

3. ConfigurationPlan
    - action_id ~ back ref for action for which plan is created
    - 




Examples:


-- > RateLimiterCluster:
<SERVICE_ID>:<ACTION>:<CONFIG_ID>:<ALG>:<UID>:<OBJ_ID>
disscord:sendMessage:default:<UUID_OF_USER>:<GROUP_ID>
disscord:sendMessage:nitro:<UUID_OF_USER>:<GROUP_ID>


identity:login:127.0.0.0:
identity:signUp:127.0.0.0:
identity:changePassword:127.0.0.0:
miniurl:createUrl:SILVER:ajkd-2eni-2dhn-jadu:
miniurl:createUrl:SILVER:ajkd-2eni-2dhn-jadu:


rate_limit:
    unit: minutes/hours/days/months
    requests_per_unit: 1/100/1000


Headers
X-RateLimit-Limit
X-RateLimit-Remaining
X-RateLimit-Reset-After
X-RateLimit-Reset


# Example Flow
User from IP 127.0.0.1 tries to login

@api.login
def login():
    ip = request.headers.get("ip)

    # Check rate limit
    rateLimitResponse = RPC.checkRateLimit(
                                                service_name = "identity",
                                                action_name  = "login",
                                                congif_name  = "defualt",
                                                uid          = "127.0.0.1",
                                                oid          = null,
                                            )