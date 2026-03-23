#!/bin/bash

curl -iX POST 'localhost:8080/webhook/' -d '{"message": "hi"}'
