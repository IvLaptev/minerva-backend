#!/bin/sh
sleep 1
echo "I'm the script with pid $$"
for i in 1 2 3 4 5; do
        sleep 1
        echo "Still running $$"
done

echo "Stopped $$"
