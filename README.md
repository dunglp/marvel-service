# Marvel Service
Technical Assessment

## Author :

-   Dzung Le - Senior FullStack Engineer

## Prerequisites : 
-   Docker installed
-   localhost:8080, localhost:3003 is not in used

## How to run ?

-   Open _terminal_ in the same folder
-   **Run** : _make run_ to start Marvel service server using Docker at http://localhost:8080 
-   **Run** : _make serve-swagger_ to open local swagger at http://localhost:3003/docs
-   Test can be done via swagger UI 

## About this project

-   **http://localhost:8080/characters :**
    -   Get all Marvel character IDs 
    -   Caching strategies : 
        - 2 strategies were considered
        - Strategy 1 :
          - look for cache 
            - if not found, get all IDs from Marvel, cache all IDs, return
            - if found :
                - get 1 ID from Marvel, results will contain Total IDs
                - compare Marvel.Total and Cache.Total :
                    - if equal, return IDs from Cache
                    - if not equal, get all IDs from Marvel, cache all IDs, return
          - https://github.com/dunglp/marvel-service/blob/main/pkg/marvel/cache/diagrams/Marvel_Svc_Caching_.png
        - Strategy 2 : (implemented in this assessment)
          - look for cache 
            - if not found, get all IDs from Marvel, cache : 
              - Cache.IDs = Marvel.IDs
              - Cache.since = time.Now
              - return
            - if found :
              - get since = Cache.since
              - query for all IDs that has been modified from _since_ timestamp (Assumption : All newly added character is included)
              - add all new IDs to cache, Cache.since = time.Now
              - return
          - https://github.com/dunglp/marvel-service/blob/main/pkg/marvel/cache/diagrams/Marvel_Svc_Caching_2.png
    
-   **http://localhost:8080/characters/{characterId} :**
    -   Get Marvel character details 

- OpenAPI specs generated using swagger