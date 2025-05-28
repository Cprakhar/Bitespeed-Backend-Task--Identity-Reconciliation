# Bitespeed Identity Reconciliation Service

A Go backend service for identifying and linking customer contacts across multiple purchases.

## Problem Statement

FluxKart.com needs to identify and track customer identity across multiple purchases, even when customers use different email addresses and phone numbers for each order.

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your settings
3. Run `go mod tidy` to install dependencies
4. Start the server with `go run cmd/server/main.go`

## API Endpoints

- `GET /health` - Health check endpoint
- `POST /identify` - Identity reconciliation endpoint (coming soon)

## Database Schema

The service uses SQLite database with a `contacts` table for storing customer contact information.

### Contact Table Schema
- `id` - Primary key (auto-increment)
- `phone_number` - Phone number (optional)
- `email` - Email address (optional)
- `linked_id` - Foreign key to another contact (for linking)
- `link_precedence` - Either 'primary' or 'secondary'
- `created_at` - Timestamp when record was created
- `updated_at` - Timestamp when record was last updated
- `deleted_at` - Soft delete timestamp (NULL if not deleted)

## Technology Stack

- **Backend**: Go 1.21
- **Database**: SQLite
- **HTTP Framework**: Standard net/http