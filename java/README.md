[← Back to Platform Root](../../README.md)

# NativeBPM Java SDK

Official Java Client SDK for NativeBPM Cloud-Native engine. Built using the `okhttp-gson` HTTP client library.

## Installation & Building

To compile the SDK locally and install it into your local Maven repository (`~/.m2`):

```bash
mvn clean install
```

Or using Gradle:

```bash
./gradlew build
```

## Maven Coordinates

To consume the package from the GitLab Maven Package Registry, add the repository configuration to your build file:

**Maven (`pom.xml`)**:
```xml
<repositories>
    <repository>
        <id>gitlab-maven</id>
        <url>https://gitlab.com/api/v4/projects/nativebpm%2Fplatform/packages/maven</url>
    </repository>
</repositories>

<dependencies>
    <dependency>
        <groupId>com.nativebpm</groupId>
        <artifactId>nativebpm-java-client</artifactId>
        <version>1.0.0</version>
    </dependency>
</dependencies>
```

**Gradle (`build.gradle`)**:
```groovy
repositories {
    maven {
        url "https://gitlab.com/api/v4/projects/nativebpm%2Fplatform/packages/maven"
    }
}

dependencies {
    implementation 'com.nativebpm:nativebpm-java-client:1.0.0'
}
```

## Usage Example

Import the custom Java classes under the `com.nativebpm.client` package namespace:

```java
import com.nativebpm.client.ApiClient;
import com.nativebpm.client.ApiException;
import com.nativebpm.client.Configuration;
import com.nativebpm.client.auth.HttpBearerAuth;
import com.nativebpm.client.api.DefaultApi;
import com.nativebpm.client.model.ProcessDefinition;
import com.nativebpm.client.model.ProcessInstance;
import com.nativebpm.client.model.StartInstanceRequest;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Main {
    public static void main(String[] args) {
        // Initialize and configure ApiClient
        ApiClient client = Configuration.getDefaultApiClient();
        client.setBasePath("http://localhost:8080");

        // Set bearer token authentication
        HttpBearerAuth bearerAuth = (HttpBearerAuth) client.getAuthentication("bearerAuth");
        bearerAuth.setBearerToken("your-api-token");

        DefaultApi apiInstance = new DefaultApi(client);

        try {
            // 1. List deployed process definitions
            List<ProcessDefinition> definitions = apiInstance.listDefinitions();
            for (ProcessDefinition def : definitions) {
                System.out.println("Definition: " + def.getName() + " (ID: " + def.getId() + ")");
            }

            // 2. Start a process instance
            StartInstanceRequest startRequest = new StartInstanceRequest();
            startRequest.setBusinessKey("REQ-2201");
            
            Map<String, Object> variables = new HashMap<>();
            variables.put("approved", true);
            startRequest.setVariables(variables);

            ProcessInstance instance = apiInstance.startInstance("expense-process", startRequest);
            System.out.println("Started instance: " + instance.getId());

        } catch (ApiException e) {
            System.err.println("API error code: " + e.getCode());
            System.err.println("API error body: " + e.getResponseBody());
            e.printStackTrace();
        }
    }
}
```

## Developer Guide: Publishing to GitLab Package Registry

### Automated Publishing via CI/CD
This package is automatically built, versioned, and published to the GitLab Maven Package Registry whenever a git tag matching the pattern `sdk/java/v*` is pushed. For example:
```bash
git tag sdk/java/v1.0.0
git push origin sdk/java/v1.0.0
```

### Manual Publishing
To manually publish a release version:
1. Configure your local `~/.m2/settings.xml` to include GitLab repository authentication:
   ```xml
   <settings>
     <servers>
       <server>
         <id>gitlab-maven</id>
         <configuration>
           <httpHeaders>
             <property>
               <name>Private-Token</name>
               <!-- Or Deploy-Token -->
               <value>your_personal_access_or_deploy_token</value>
             </property>
           </httpHeaders>
         </configuration>
       </server>
     </servers>
   </settings>
   ```
2. Run Maven set-version and deploy:
   ```bash
   mvn versions:set -DnewVersion=1.0.0
   mvn deploy -DskipTests
   ```

