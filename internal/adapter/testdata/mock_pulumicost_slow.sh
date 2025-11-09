#!/bin/bash
# Mock slow pulumicost-core binary for testing timeouts

# Simulate slow operation
sleep 10

# Return mock result (won't be reached if context is canceled)
cat <<EOF
{
  "total_monthly": 42.50,
  "currency": "USD",
  "resources": []
}
EOF
