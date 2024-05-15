#!/bin/bash

API_ENDPOINT="http://localhost/api/v1/purchases"
DATA='{
    "order_items": [
        {
            "product_id":  "f868ff38-32e4-4527-82aa-4ed8dd2d6247",
            "quantity": 4
        },
        {
            "product_id":"3f1cb659-f12e-47a3-909e-117fe1d0db79",
            "quantity": 10
        }
    ],
    "payment": {
        "currency_code": "VND"
    }
}'

for ((i=1; i<=1000; i++)); do
    curl --location --request POST "$API_ENDPOINT" \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer v2.local.vlOD1v-3B7C8bXU8L-tfvI1IOrdyp7cjz0FhKhWVKV2n4cF78gBKsuXsQyrLGa4xlPF9rkuI2g0DqyCukKil89qzJWLJSJsKyDrZWTPdzDHoYtc4Y8bo_3zE2NTrQo60IwEE5laZdZmqmAOPafEa1NQ8Ww4D65WTIye-7iLJWfCFo0bxsm2zWMY4_YCaWoh_pD89HrYTeHYTHx1ISr_r4hB6xW_c7gWOMjsQ9SM27S_TaBR0jjFsMMPBnRiJv5wD24U86vS1KTSy2p8sOCWn_t9E5cue3B1EgkTQ4ZpOHz99.bnVsbA' \
    --data-raw "$DATA" &
done

wait