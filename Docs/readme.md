# **System Design Document: Distributed Rate Limiter**


## **Problem Statement**
- Design and implement a **distributed rate limiter** at the application layer that can serve any number of services. 
- The rate limiter should enforce quotas and prevent abuse by limiting the number of requests per client or service within a specified time window. 
- It must be scalable, efficient, and resilient to failures.

---

# Functional Requirements
1. Enforce rate limits at the **SERVICE_ID:ACTION:CONFIG_ID:ALG:UID:OBJ_ID**.
2. Support multiple types of rate-limiting strategies:
    - **Fixed Window**: Limit requests in discrete intervals.
    - **Sliding Window**: Smoothen limits over time to avoid burst traffic.
3. Provide rate-limiting headers:
  - `X-RateLimit-Limit`: Maximum allowed requests.
  - `X-RateLimit-Remaining`: Remaining requests in the current window.
  - `X-RateLimit-Reset`: Time when the limit resets.
4. The rate limit configurations needs to be cutsomisable as deep as possible.



# Non-Functional Requirements
    1. Performance
        - Low latency - Since each request is rate limited we need to keep overall latency of the service request as low as possible (response times under 50 ms).
        - High throughput to handle concurrent requests.
    2. Scalability
        - Scale horizontally to support increasing numbers of requests .
        - Ensure the system can handle peak loads during high-traffic events.     

# Good to have requirements
    1. Provide APIs to retrieve usage metrics, like the current usage or remaining quota per service.
    2. Push the rate limit pattern to Analytics team to enhance the overall user experience of the platform.
    3. Additionally support rate limiting in gateway level.



# Configuration Example
```
services:
- name: identity
  id: 1
  actions:
  - name: login
    id: 1
    rate_limit:
    - name: default
      unit: minutes
      unit_multiplier: 2
      request_per_unit: 3
      algorithm: fixed_window

    - name: basic
      unit: minutes
      unit_multiplier: 1
      request_per_unit: 10
      algorithm: sliding_window

    - name: premium
      unit: hours
      unit_multiplier: 1
      request_per_unit: 20
      algorithm: fixed_window

  - name: disscord
    id: 2
    actions:
    - name: sendMessage
      id: 1
      rate_limit:
      - name: default
        unit: minutes
        unit_multiplier: 1
        request_per_unit: 60
        algorithm: fixed_window

      - name: nitro
        unit: minutes
        unit_multiplier: 1
        request_per_unit: 120
        algorithm: sliding_window

```


# High Level Design
    1. Storage Layer
        - Use Redis as primary database, Since the qureying pattern is only Key-value based and we need low latency data access we shall use Redis. 
    2. Transport Layer
        - Expose gRPC endpoint for checking rate limits
        - Since it is going to be a service to service intercommunication,  gRPC is the best candidate for this requirement. 
        

# Tech Stack Used
1. Golang with gRPC 
2. Redis
3. Docker
4. Kubernetes


# API Design
```
gRPC 
1. CheckRateLimit
    - Request
        - service_name             
        - action_name  
        - config_name (For customised rate limiting )      
        - unique_identifier (Say IP address or UUID of the user)                
        - object_identifier     
    - Response
        - is_rate_limited
        - limit
        - remaining
        - reset_at
        - reset_after
        
```


