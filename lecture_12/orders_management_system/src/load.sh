#!/bin/sh
for (( i=1; i <= 10000; i++ ))
do
grpc_cli call --json_input --json_output localhost:8082 github.com.moguchev.microservices.orders_management_system.OrdersManagementSystemService/CreateOrder '{"delivery_info":{"delivery_date":{"nanos":0,"seconds":"1717020669"},"delivery_variant_id":"1"},"items":[{"warehouse_id":"1","quantity":1,"id":"1"},{"quantity":2,"id":"2","warehouse_id":"2"},{"id":"3","warehouse_id":"3","quantity":3}],"user_id":"9"}'
done
