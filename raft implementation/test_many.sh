#!/bin/bash

# File to store the log output
LOGFILE="logfile.log"

# Clear previous log file (optional, remove if you want to append)
> "$LOGFILE"

# Loop from 1 to 100
for i in {1..100}
do
    # Log which run it is
    echo "Running test iteration $i" | tee -a "$LOGFILE"

    # Run the test and log the time it takes
    { time go test ; } 2>&1 | tee -a "$LOGFILE"

    # Check if the test failed
    if grep -q "FAIL" "$LOGFILE"; then
        echo "Test failed in iteration $i. Stopping the process."
        exit 1
    fi

    # Separate each run's log output with a line (optional)
    echo "----------------------------------------" | tee -a "$LOGFILE"
done

echo "All test runs completed. Check $LOGFILE for details."
