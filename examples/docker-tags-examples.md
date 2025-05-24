# Docker Image Tagging Examples

This document explains how Docker images are tagged in different scenarios.

## Main Branch Commits

When pushing to the main branch:
- `ghcr.io/rknightion/adsb2loki:main` - branch tag
- `ghcr.io/rknightion/adsb2loki:sha-abc1234` - commit SHA tag

## Pull Requests

When creating a pull request:
- `ghcr.io/rknightion/adsb2loki:pr-123` - PR number tag

## Version Tags/Releases

When creating a version tag (e.g., `v1.2.3`):
- `ghcr.io/rknightion/adsb2loki:1.2.3` - full version
- `ghcr.io/rknightion/adsb2loki:1.2` - major.minor
- `ghcr.io/rknightion/adsb2loki:1` - major only
- `ghcr.io/rknightion/adsb2loki:latest` - latest tag (only for releases)

Note: SHA tags are NOT created for version releases to avoid the invalid tag format issue. 