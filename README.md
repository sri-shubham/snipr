# Snipr
## URL Shortener

[![Go Report Card](https://goreportcard.com/badge/github.com/sri-shubham/neon)](https://goreportcard.com/report/github.com/sri-shubham/snipr)
[![GitHub issues](https://img.shields.io/github/issues/sri-shubham/neon)](https://github.com/sri-shubham/snipr/issues)
[![GitHub stars](https://img.shields.io/github/stars/sri-shubham/neon)](https://github.com/sri-shubham/snipr/stargazers)


This is a simple URL shortening service.

## Features:
- Modular and Low Coupling: The application is designed with a modular architecture, making it easy to replace any of the components (e.g., storage, caching) with your own implementation.
- Individually Testable Packages: Each package in the application is designed to be independently testable, ensuring better maintainability and reliability.
- SHA-256 Hash for URL Shortening: Snipr uses the SHA-256 hashing algorithm to generate short URLs, providing a secure and efficient URL shortening mechanism.
Customizable: Snipr allows you to customize various aspects of the URL shortening service, such as the domain, URL format, and more.

## Getting Started:
To get started with Snipr, follow these steps:
- Clone the repository: git clone https://github.com/sri-shubham/snipr.git
- Navigate to the project directory: cd snipr
- `docker-compose up -d`

using current config it starts up service and postgres containers. There is redis storage interface implemented as well which can be drop in replacement for postgres storage. Although reporting is not yet implemented for redis I will come around to implement a complete drop in replacement redis available.
