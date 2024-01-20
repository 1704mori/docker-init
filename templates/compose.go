package templates

const COMPOSE = `
services:
  server:
    build:
      context: .
      target: final
    ports:
      - {{.Port}}:{{.Port}}
`
