# Run each command separateley. The GSIs take some time to finish creating. You can inspect their status in the AWS Dynamo UI
# Export AWS credentials with Dynamo write permissions and also the TABLE_PREFIX env variable
# Create the GSIs before running the backfill binary

aws dynamodb update-table \
    --table-name $TABLE_PREFIX-audit-logs \
    --attribute-definitions AttributeName=timestamp,AttributeType=S AttributeName=user#org,AttributeType=S \
    --global-secondary-index-updates \
        "[
            {
                \"Create\": {
                    \"IndexName\": \"user-org-index\",
                    \"KeySchema\": [{\"AttributeName\":\"user#org\",\"KeyType\":\"HASH\"},
                                    {\"AttributeName\":\"timestamp\",\"KeyType\":\"RANGE\"}],
                    \"Projection\":{
                        \"ProjectionType\":\"ALL\"
                    }
                }
            }
        ]"

aws dynamodb update-table \
    --table-name $TABLE_PREFIX-audit-logs \
    --attribute-definitions AttributeName=timestamp,AttributeType=S AttributeName=rack-name#org,AttributeType=S \
    --global-secondary-index-updates \
        "[
            {
                \"Create\": {
                    \"IndexName\": \"rack-name-org-index\",
                    \"KeySchema\": [{\"AttributeName\":\"rack-name#org\",\"KeyType\":\"HASH\"},
                                    {\"AttributeName\":\"timestamp\",\"KeyType\":\"RANGE\"}],
                    \"Projection\":{
                        \"ProjectionType\":\"ALL\"
                    }
                }
            }
        ]"


aws dynamodb update-table \
    --table-name $TABLE_PREFIX-audit-logs \
    --attribute-definitions AttributeName=timestamp,AttributeType=S AttributeName=user#rack-name#org,AttributeType=S \
    --global-secondary-index-updates \
        "[
            {
                \"Create\": {
                    \"IndexName\": \"user-rack-name-org-index\",
                    \"KeySchema\": [{\"AttributeName\":\"user#rack-name#org\",\"KeyType\":\"HASH\"},
                                    {\"AttributeName\":\"timestamp\",\"KeyType\":\"RANGE\"}],
                    \"Projection\":{
                        \"ProjectionType\":\"ALL\"
                    }
                }
            }
        ]"

