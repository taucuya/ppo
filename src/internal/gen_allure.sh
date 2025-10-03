#!/bin/bash

gotestsum --junitfile allure-results/junit.xml --format standard-verbose ./...

allure generate allure-results --output allure-report --clean

allure open allure-report