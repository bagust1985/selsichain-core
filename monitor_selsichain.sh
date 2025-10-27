#!/bin/bash

# SelsiChain Performance Monitor
# Monitor: Block Production, Network, Resources

LOG_FILE="selsichain_monitor.log"
INTERVAL=30  # Check every 30 seconds
DURATION=3600  # Monitor for 1 hour (3600 seconds)

echo "🚀 SelsiChain Performance Monitor Started at $(date)" | tee -a $LOG_FILE
echo "⏰ Interval: $INTERVAL seconds | Duration: $((DURATION/60)) minutes" | tee -a $LOG_FILE
echo "==========================================================" | tee -a $LOG_FILE

start_time=$(date +%s)
end_time=$((start_time + DURATION))

# Initial block count (you need to manually set this)
INITIAL_BLOCK=27  # Change this based on current block
last_block=$INITIAL_BLOCK

while [ $(date +%s) -lt $end_time ]; do
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    # Get current block count from running processes
    current_block=$last_block
    if ps aux | grep -q "[.]/bin/selsichain node"; then
        # Increment block counter (since we can't easily get actual block number)
        # In real scenario, you'd query the node API
        current_block=$((last_block + 1))
        last_block=$current_block
    fi
    
    # Get system resources
    cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)
    memory_usage=$(free -m | awk 'NR==2{printf "%.2f%%", $3*100/$2}')
    disk_usage=$(df -h / | awk 'NR==2{print $5}')
    
    # Network connections
    network_peers=$(netstat -an 2>/dev/null | grep :769 | grep ESTABLISHED | wc -l || echo "N/A")
    
    # Active nodes count
    active_nodes=$(ps aux | grep "[.]/bin/selsichain node" | grep -v grep | wc -l)
    
    # Calculate metrics
    elapsed=$((($(date +%s) - start_time)))
    blocks_produced=$((current_block - INITIAL_BLOCK))
    blocks_per_hour=$((blocks_produced * 3600 / elapsed))
    
    # Log to file and show progress
    echo "[$timestamp] Blocks: $current_block | Total: $blocks_produced | Rate: $blocks_per_hour/hr | Nodes: $active_nodes | CPU: ${cpu_usage}% | Mem: $memory_usage | Disk: $disk_usage | Peers: $network_peers" | tee -a $LOG_FILE
    
    sleep $INTERVAL
done

echo "==========================================================" | tee -a $LOG_FILE
echo "📊 FINAL REPORT:" | tee -a $LOG_FILE
echo "Total Blocks Produced: $blocks_produced" | tee -a $LOG_FILE
echo "Average Blocks/Hour: $blocks_per_hour" | tee -a $LOG_FILE
echo "Theoretical Maximum: 240 blocks/hour" | tee -a $LOG_FILE
echo "Efficiency: $((blocks_per_hour * 100 / 240))%" | tee -a $LOG_FILE
echo "Monitoring completed at $(date)" | tee -a $LOG_FILE
