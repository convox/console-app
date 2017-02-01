#!/bin/sh

die() {
  echo $*
  exit 1
}

stack=$1

[ -z $stack ] && die "usage: resource-environmnent.sh <stack-name>"
\
aws cloudformation describe-stacks --stack-name $stack | jq -r '.Stacks[0].Outputs[] | [(.OutputKey | gsub("(?<pre>.)(?<up>[A-Z])";"\(.pre)_\(.up)") | ascii_upcase), .OutputValue] | join("=")'
