package main

import (
	"fmt"
	"os"
	"strings"
)

var tableTemplate = `    "ConsolePrivate%%TABLE_CAPS%%ScalableTargetRead": {
      "Type": "AWS::ApplicationAutoScaling::ScalableTarget",
      "Properties": {
        "MaxCapacity": 50,
        "MinCapacity": 5,
        "ResourceId": { "Fn::Sub": "table/${ConsolePrivate%%TABLE_CAPS%%DynamoDBTable}" },
        "RoleARN": { "Fn::GetAtt": [ "DynamoScalingRole", "Arn" ] },
        "ScalableDimension": "dynamodb:table:ReadCapacityUnits",
        "ServiceNamespace": "dynamodb"
      }
    },
    "ConsolePrivate%%TABLE_CAPS%%ScalableTargetWrite": {
      "Type": "AWS::ApplicationAutoScaling::ScalableTarget",
      "Properties": {
        "MaxCapacity": 50,
        "MinCapacity": 2,
        "ResourceId": { "Fn::Sub": "table/${ConsolePrivate%%TABLE_CAPS%%DynamoDBTable}" },
        "RoleARN": { "Fn::GetAtt": [ "DynamoScalingRole", "Arn" ] },
        "ScalableDimension": "dynamodb:table:WriteCapacityUnits",
        "ServiceNamespace": "dynamodb"
      }
    },
    "ConsolePrivate%%TABLE_CAPS%%ScalingRead": {
      "Type": "AWS::ApplicationAutoScaling::ScalingPolicy",
      "Properties": {
        "PolicyName": { "Fn::Sub": "${TablePrefix}-%%TABLE_NORMAL%%-scaling-read" },
        "PolicyType": "TargetTrackingScaling",
        "ScalingTargetId": { "Ref": "ConsolePrivate%%TABLE_CAPS%%ScalableTargetRead" },
        "TargetTrackingScalingPolicyConfiguration": {
          "TargetValue": 70.0,
          "ScaleInCooldown": 5,
          "ScaleOutCooldown": 60,
          "PredefinedMetricSpecification": { "PredefinedMetricType": "DynamoDBReadCapacityUtilization" }
        }
      }
    },
    "ConsolePrivate%%TABLE_CAPS%%ScalingWrite": {
      "Type": "AWS::ApplicationAutoScaling::ScalingPolicy",
      "Properties": {
        "PolicyName": { "Fn::Sub": "${TablePrefix}-%%TABLE_NORMAL%%-scaling-write" },
        "PolicyType": "TargetTrackingScaling",
        "ScalingTargetId": { "Ref": "ConsolePrivate%%TABLE_CAPS%%ScalableTargetWrite" },
        "TargetTrackingScalingPolicyConfiguration": {
          "TargetValue": 70.0,
          "ScaleInCooldown": 5,
          "ScaleOutCooldown": 60,
          "PredefinedMetricSpecification": { "PredefinedMetricType": "DynamoDBWriteCapacityUtilization" }
        }
      }
    },
`

var indexTemplate = `    "ConsolePrivate%%TABLE_CAPS%%Index%%INDEX_CAPS%%ScalableTargetRead": {
      "Type": "AWS::ApplicationAutoScaling::ScalableTarget",
      "Properties": {
        "MaxCapacity": 50,
        "MinCapacity": 5,
        "ResourceId": { "Fn::Sub": "table/${ConsolePrivate%%TABLE_CAPS%%DynamoDBTable}/index/%%INDEX_NORMAL%%-index" },
        "RoleARN": { "Fn::GetAtt": [ "DynamoScalingRole", "Arn" ] },
        "ScalableDimension": "dynamodb:index:ReadCapacityUnits",
        "ServiceNamespace": "dynamodb"
      }
    },
    "ConsolePrivate%%TABLE_CAPS%%Index%%INDEX_CAPS%%ScalableTargetWrite": {
      "Type": "AWS::ApplicationAutoScaling::ScalableTarget",
      "Properties": {
        "MaxCapacity": 50,
        "MinCapacity": 2,
        "ResourceId": { "Fn::Sub": "table/${ConsolePrivate%%TABLE_CAPS%%DynamoDBTable}/index/%%INDEX_NORMAL%%-index" },
        "RoleARN": { "Fn::GetAtt": [ "DynamoScalingRole", "Arn" ] },
        "ScalableDimension": "dynamodb:index:WriteCapacityUnits",
        "ServiceNamespace": "dynamodb"
      }
    },
    "ConsolePrivate%%TABLE_CAPS%%Index%%INDEX_CAPS%%ScalingRead": {
      "Type": "AWS::ApplicationAutoScaling::ScalingPolicy",
      "Properties": {
        "PolicyName": { "Fn::Sub": "${TablePrefix}-%%TABLE_NORMAL%%-index-%%INDEX_NORMAL%%-scaling-read" },
        "PolicyType": "TargetTrackingScaling",
        "ScalingTargetId": { "Ref": "ConsolePrivate%%TABLE_CAPS%%Index%%INDEX_CAPS%%ScalableTargetRead" },
        "TargetTrackingScalingPolicyConfiguration": {
          "TargetValue": 70.0,
          "ScaleInCooldown": 5,
          "ScaleOutCooldown": 60,
          "PredefinedMetricSpecification": { "PredefinedMetricType": "DynamoDBReadCapacityUtilization" }
        }
      }
    },
    "ConsolePrivate%%TABLE_CAPS%%Index%%INDEX_CAPS%%ScalingWrite": {
      "Type": "AWS::ApplicationAutoScaling::ScalingPolicy",
      "Properties": {
        "PolicyName": { "Fn::Sub": "${TablePrefix}-%%TABLE_NORMAL%%-index-%%INDEX_NORMAL%%-scaling-write" },
        "PolicyType": "TargetTrackingScaling",
        "ScalingTargetId": { "Ref": "ConsolePrivate%%TABLE_CAPS%%Index%%INDEX_CAPS%%ScalableTargetWrite" },
        "TargetTrackingScalingPolicyConfiguration": {
          "TargetValue": 70.0,
          "ScaleInCooldown": 5,
          "ScaleOutCooldown": 60,
          "PredefinedMetricSpecification": { "PredefinedMetricType": "DynamoDBWriteCapacityUtilization" }
        }
      }
    },
`

func main() {
	switch len(os.Args) {
	case 2:
		table(os.Args[1])
	case 3:
		index(os.Args[1], os.Args[2])
	}
}

func table(t string) {
	out := tableTemplate

	out = strings.Replace(out, "%%TABLE_NORMAL%%", t, -1)
	out = strings.Replace(out, "%%TABLE_CAPS%%", upperName(t), -1)

	fmt.Print(out)
}

func index(t, i string) {
	out := indexTemplate

	out = strings.Replace(out, "%%TABLE_NORMAL%%", t, -1)
	out = strings.Replace(out, "%%TABLE_CAPS%%", upperName(t), -1)
	out = strings.Replace(out, "%%INDEX_NORMAL%%", i, -1)
	out = strings.Replace(out, "%%INDEX_CAPS%%", upperName(i), -1)

	fmt.Print(out)
}

func upperName(name string) string {
	if name == "" {
		return ""
	}

	// myapp -> Myapp; my-app -> MyApp
	us := strings.ToUpper(name[0:1]) + name[1:]

	for {
		i := strings.Index(us, "-")

		if i == -1 {
			break
		}

		s := us[0:i]

		if len(us) > i+1 {
			s += strings.ToUpper(us[i+1 : i+2])
		}

		if len(us) > i+2 {
			s += us[i+2:]
		}

		us = s
	}

	return us
}
