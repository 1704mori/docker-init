package templates

const NODE = `
# syntax=docker/dockerfile:1

ARG NODE_VERSION={{.LanguageVersion}}

FROM node:${NODE_VERSION}-alpine AS builder

WORKDIR /usr/src/app

# Copy only the package files for detection
COPY package.json yarn.lock pnpm-lock.yaml ./

# Detect the lock file and install dependencies
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=yarn.lock,target=yarn.lock \
    --mount=type=bind,source=pnpm-lock.yaml,target=pnpm-lock.yaml \
    if [ -f yarn.lock ]; then \
        yarn install --frozen-lockfile; \
    elif [ -f pnpm-lock.yaml ]; then \
        pnpm install; \
    else \
        npm ci --omit=dev; \
    fi

# Use the appropriate base image for the production stage
FROM node:${NODE_VERSION}-alpine

ENV NODE_ENV production

WORKDIR /usr/src/app

# Copy only the necessary files from the build stage
COPY --from=builder /usr/src/app/node_modules ./node_modules
COPY . .

USER node

EXPOSE {{.Port}}

CMD node index.js
`
