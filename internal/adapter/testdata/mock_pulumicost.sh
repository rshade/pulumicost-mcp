#!/bin/bash
# Mock pulumicost-core binary for testing

# Read input from stdin if provided
INPUT=$(cat)

# Return mock cost analysis result
cat <<EOF
{
  "total_monthly": 42.50,
  "total_hourly": 0.058,
  "currency": "USD",
  "resources": [
    {
      "urn": "urn:pulumi:dev::myapp::aws:ec2/instance:Instance::web-server",
      "name": "web-server",
      "type": "aws:ec2/instance:Instance",
      "provider": "aws",
      "monthly_cost": 10.50,
      "hourly_cost": 0.014,
      "region": "us-east-1",
      "tags": {
        "environment": "dev",
        "team": "platform"
      }
    },
    {
      "urn": "urn:pulumi:dev::myapp::aws:rds/instance:Instance::db",
      "name": "db",
      "type": "aws:rds/instance:Instance",
      "provider": "aws",
      "monthly_cost": 32.00,
      "hourly_cost": 0.044,
      "region": "us-east-1",
      "tags": {
        "environment": "dev",
        "team": "backend"
      }
    }
  ],
  "breakdown": {
    "daily": [
      {"date": "2024-01-01", "amount": 1.42},
      {"date": "2024-01-02", "amount": 1.42}
    ]
  }
}
EOF
