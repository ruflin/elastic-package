// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package install

const snapshotYml = `version: '2.3'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:${STACK_VERSION}
    healthcheck:
      test: ["CMD", "curl", "-f", "-u", "elastic:changeme", "http://127.0.0.1:9200/"]
      retries: 300
      interval: 1s
    environment:
    - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
    - "network.host="
    - "transport.host=127.0.0.1"
    - "http.host=0.0.0.0"
    - "indices.id_field_data.enabled=true"
    - "xpack.license.self_generated.type=trial"
    - "xpack.security.enabled=true"
    - "xpack.security.authc.api_key.enabled=true"
    - "ELASTIC_PASSWORD=changeme"
    - "script.context.template.max_compilations_rate=unlimited"
    - "script.context.ingest.cache_max_size=2000"
    - "script.context.processor_conditional.cache_max_size=2000"
    - "script.context.template.cache_max_size=2000"
    ports:
      - "127.0.0.1:9200:9200"

  elasticsearch_is_ready:
    image: tianon/true
    depends_on:
      elasticsearch:
        condition: service_healthy

  kibana:
    image: docker.elastic.co/kibana/kibana:${STACK_VERSION}
    depends_on:
      elasticsearch:
        condition: service_healthy
      package-registry:
        condition: service_healthy
    healthcheck:
      test: "curl -f http://127.0.0.1:5601/login | grep kbn-injected-metadata 2>&1 >/dev/null"
      retries: 600
      interval: 1s
    volumes:
      - ./kibana.config.yml:/usr/share/kibana/config/kibana.yml
    ports:
      - "127.0.0.1:5601:5601"

  kibana_is_ready:
    image: tianon/true
    depends_on:
      kibana:
        condition: service_healthy

  package-registry:
    build:
      context: .
      dockerfile: Dockerfile.package-registry
    healthcheck:
      test: ["CMD", "curl", "-f", "http://127.0.0.1:8080"]
      retries: 300
      interval: 1s
    ports:
      - "127.0.0.1:8080:8080"

  package-registry_is_ready:
    image: tianon/true
    depends_on:
      package-registry:
        condition: service_healthy

  elastic-agent:
    image: docker.elastic.co/beats/elastic-agent:${STACK_VERSION}
    depends_on:
      elasticsearch:
        condition: service_healthy
      kibana:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://127.0.0.1:8080"]
      retries: 300
      interval: 1s
    #  test: "sh -c 'grep \"Agent is starting\" /usr/share/elastic-agent/elastic-agent.log*'"ff
    #  retries: 30
    #  interval: 1s
    hostname: docker-fleet-agent
    environment:
# Reference: https://github.com/elastic/beats/blob/master/x-pack/elastic-agent/pkg/agent/cmd/container.go#L55
# Generic configss
    - "ELASTICSEARCH_HOST=http://elasticsearch:9200"
    - "ELASTICSEARCH_USERNAME=elastic"
    - "ELASTICSEARCH_PASSWORD=changeme"
    - "KIBANA_HOST=http://kibana:5601"
    - "KIBANA_USERNAME=elastic"
    - "KIBANA_PASSWORD=changeme"

# Setup of Fleet Integrations in Kibana
    - "KIBANA_FLEET_SETUP=1"

# Fleet Server enabling
    - "FLEET_SERVER_ENABLE=1"
    - "FLEET_URL=http://localhost:8210"
#- "FLEET_SERVER_INSECURE_HTTP=1"

# Enroll the Agent into itself
    #- "FLEET_ENROLL=1"
    #- "FLEET_ENROLL_INSECURE=1"marcin
    #- "FLEET_SETUP=1"
    #- "FLEET_URL=http://kibana:5601"
    #- "KIBANA_HOST=http://kibana:5601"
    #- "ELASTICSEARCH_HOST=http://elasticsearch:9200"
    #- "FLEET_SERVER_ELASTICSEARCH_USERNAME=elastic"
    #- "FLEET_SERVER_ELASTICSEARCH_PASSWORD=changeme"

    volumes:
    - type: bind
      source: ../tmp/service_logs/
      target: /tmp/service_logs/

  elastic-agent_is_ready:
    image: tianon/true
    depends_on:
      elastic-agent:
        condition: service_healthy
`
