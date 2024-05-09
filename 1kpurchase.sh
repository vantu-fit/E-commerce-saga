#!/bin/bash

API_ENDPOINT="http://localhost/api/v1/purchases"
DATA='{
    "order_items": [
        {
            "product_id": "7b8f5a2d-f7c1-4229-86ca-25bb8384c243",
            "quantity": 4
        },
        {
            "product_id": "60d44700-4693-4a59-aa5b-57994f200129",
            "quantity": 10
        }
    ],
    "payment": {
        "currency_code": "VND"
    }
}'

for ((i=1; i<=100; i++)); do
    curl --location --request POST "$API_ENDPOINT" \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer v2.local.MLuMunbiIRsmqQ5EZ_EZCivYu5AKbbYUDvNH__ev4eCdUDvg0mSdIQc60KNsck3KSQmRexJDtvvbEqqx_SDKPo2PTlCqUqQaKmv5weCHl7L-s5iwbeDOXeyKH0Fx31PZ_UcF2qR7YisY88fw9C0SFfPTRDRspuyS4NtE8o5_6N-dK5pmVn11Rf6pZ29KB_Lyg1GDgUiYxq7SQzvluS0uKm1iu00EW-lDFSA-mUKzmU0HoPLqP81ZJqiDYsjJ60wt3LX4fr1AUXSYLnm8VsqqPg.bnVsbA' \
    --data-raw "$DATA" &
done

wait