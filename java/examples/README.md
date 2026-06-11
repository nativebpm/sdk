[← Back to Java SDK](../README.md) | [← Back to Platform Root](../../../README.md)

# NativeBPM Java SDK Example

This directory contains a complete, runnable Java console application that demonstrates how to interact with the NativeBPM Engine API using the generated Java client SDK.

## Prerequisites
- **Java JDK 17** or higher
- **Maven 3.6** or higher
- Running **NativeBPM Server** on `http://localhost:8080` (with API token `test-bearer-token`)

## Setup & Running

1. **Install the SDK into your local Maven cache (`~/.m2`)**:
   Before compiling the example, build and install the generated Java SDK client:
   ```bash
   cd ../
   mvn clean install
   ```

2. **Run the Example**:
   Navigate back to this folder and execute the Maven command:
   ```bash
   cd examples
   mvn compile exec:java
   ```

The script will automatically deploy a simple process schema, start an instance, and print active definitions/instances retrieved from the server.
